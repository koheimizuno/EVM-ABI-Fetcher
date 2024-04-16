package fetch

import (
	myCache "code/src/cache"
	myDB "code/src/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/petermattis/goid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ApiResponse
// @dev For parse the data from Etherscan
type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// FetcherCli
// @dev Config for fetching the data from Etherscan
type FetcherCli struct {
	ApiKey string
	RpcUrl string
	mu     sync.RWMutex
}

var cache = *myCache.NewABICache()
var log = logrus.New()

// GetABIAtStartOfBlock
// @notice We do not use the block parameter for now
// @param db The SQLite3's handle
// @param chainID The chain ID
// @param contractAddress The contract whose ABI you want to get
// @param block Get the ABI in a certain block height.
func (f *FetcherCli) GetABIAtStartOfBlock(db *gorm.DB, chainID int, contractAddress common.Address) ([]byte, error) {

	// [1. In memory]
	abi, isFound := cache.Get(chainID, contractAddress)
	if isFound {
		log.Info("Thread ", goid.Get(), ". Found ABI in cache, length:", len(abi))
		return abi, nil
	}

	// [2. In DB] Check if the ABI exists in the database for the given chainID, contract address, and function signature.
	//               If found, return the ABI from the database

	addressBytes := contractAddress.Bytes() // Convert Ethereum addresses to a format that can be queried by the database

	var deployment myDB.ContractDeployment
	// Query in the database
	result := db.Where("chain_id = ? AND contract_address = ?", chainID, addressBytes).First(&deployment)
	if result.Error == nil {
		// Search for the specific ABI by ContractDeployment.FunctionSignatureID
		var signature myDB.FunctionSignature
		// In some cases, there may be multiple identical pieces of data, and we will only take one of them.
		result = db.Where("id = ?", deployment.FunctionSignatureID).First(&signature)
		if result.Error == nil { // result.Error is equal to nil means that there is an ABI in the database of the input contractAddress
			///////////////////// Write Lock //////////////////////////
			f.mu.Lock()
			defer f.mu.Unlock()

			// Second check to prevent this situation: Multiple threads complete read operations simultaneously and then queue up for write operations,
			// which can cause duplicate cache writes
			abiInSecondCheck, isFoundInSecondCheck := cache.Get(chainID, contractAddress)
			if isFoundInSecondCheck {
				log.Info("Thread ", goid.Get(), ". Found ABI in cache, length:", len(abiInSecondCheck))

				return abiInSecondCheck, nil
			} else {
				// Set to cache
				cache.Set(chainID, contractAddress, []byte(signature.FunctionABI))
			}

			///////////////////// Write Lock //////////////////////////

			log.Info("Thread ", goid.Get(), ". Found ABI in DB, length:", len(signature.FunctionABI))
			return []byte(signature.FunctionABI), nil
		}
	}

	// [3. Etherscan] If the ABi is not found in the database, query Etherscan to retrieve the ABI
	//              If Etherscan returns the ABI, store it in the database and return it

	// Create a new record of ContractBytecode
	bytecode, err := queryRuntimeCode(f.RpcUrl, contractAddress)
	if string(bytecode) == "" {
		log.Error("Thread ", goid.Get(), ".The address is an EOA")
		return nil, errors.Wrap(errors.New("Check the address you input"), "Not contract")
	}

	// If tht contract has not been verified, it will return: Contract source code not verified
	abi, err = queryABIFromEtherscan(f.ApiKey, chainID, contractAddress)
	if string(abi) == "Contract source code not verified" {
		log.Error("Thread ", goid.Get(), ".The contract has not been verified")
		return nil, errors.Wrap(errors.New("Not verify"), "Not verify")
	}

	// If Etherscan does not have the ABI, return an appropriate error
	if err != nil {
		log.Error("Thread ", goid.Get(), ". Fail to fetch ABI by Etherscan API. ChainID:", chainID, "contractAddress:", contractAddress, "API KEY:", f.ApiKey, "RPC URL:", f.RpcUrl)
		return nil, errors.Wrap(errors.New("Fail to fetch ABI by Etherscan API"), "Fetch fail")
	}

	log.Info("Thread ", goid.Get(), ". Fetch ABI in Etherscan, length:", len(abi))

	// Prepare some UUID
	functionSignatureID := uuid.New()
	contractBytecodeID := uuid.New()

	if err != nil {
		log.Error("Thread ", goid.Get(), ". Fail to fetch Bytecode by Etherscan API. ChainID:", chainID, "contractAddress:", contractAddress, "API KEY:", f.ApiKey, "RPC URL:", f.RpcUrl)
		return nil, errors.Wrap(errors.New("Fail to fetch Bytecode by Etherscan API"), "Fetch fail")
	}

	///////////////////// Write Lock //////////////////////////
	// No secondary check( Whether the database contains the data or not) was performed on the database during the write operation,
	// so there may be multiple identical data entries in the database.
	// But this has no impact, because when we read the database, we only took one of them.
	// This situation will not occur multiple times since only has a small probability of occurring when writing to the database for the first time.

	f.mu.Lock()
	defer f.mu.Unlock()

	newContractBytecode := myDB.ContractBytecode{
		ID:       contractBytecodeID,
		Bytecode: bytecode,
	}
	db.Create(&newContractBytecode)

	// Store the ABI that fetched from Etherscan into the database
	newSignature := myDB.FunctionSignature{
		ID:          functionSignatureID,
		Signature:   nil,         // TODO How to deal with this， what to store: single signature or a array of signatures
		FunctionABI: string(abi), // We store the whole ABI of the contract
	}
	db.Create(&newSignature)

	// Create a new record of ContractDeployment
	newDeployment := myDB.ContractDeployment{
		ChainID:             chainID,
		ContractAddress:     addressBytes,
		ContractBytecodeID:  contractBytecodeID,
		FunctionSignatureID: functionSignatureID,
	}
	db.Create(&newDeployment)

	abiInSecondCheck, isFound := cache.Get(chainID, contractAddress)
	if isFound {
		log.Info("Thread ", goid.Get(), ". Found ABI in cache, length:", len(abiInSecondCheck))
		return abiInSecondCheck, nil
	} else {
		cache.Set(chainID, contractAddress, abi)
	}
	///////////////////// Write Lock //////////////////////////

	return abi, nil
}

// @dev Check ChainID and get the format the request url
// @notice Only support Ethereum network now
// @notice Sometimes we could fetch data in Etherscan without an API KEY
func checkChainIDAndGetReqURL(apiKey string, chainID int, contractAddress common.Address) (string, error) {
	var requestURL string

	if apiKey == "" {
		log.Warning("The request may be fail without an API KEY")
	}

	if chainID == 1 { // Ethereum
		requestURL = fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", contractAddress, apiKey)
	} else if chainID == 56 { // BSC
		// TODO
	} else if chainID == 42161 { // Arbitrum
		// TODO
	} else if chainID == 8453 { // Base
		// TODO
	} else if chainID == 43114 { // Avalanche
		// TODO
	} else if chainID == 137 { // Polygon
		// TODO
	} else {
		log.Error("Invalid chainID or API KEY. chainID:", chainID, "API KEY:", apiKey, "contractAddress", contractAddress)
		return "", errors.Wrap(errors.New("Check the input"), "Fail to checkChainIDAndGetReqURL")
	}

	return requestURL, nil
}

// @dev Query a contract's ABI
func queryABIFromEtherscan(apiKey string, chainID int, contractAddress common.Address) ([]byte, error) {
	requestURL, err := checkChainIDAndGetReqURL(apiKey, chainID, contractAddress)
	if err != nil {
		log.Error("Invalid requestURL:", requestURL)
		return []byte{}, errors.Wrap(errors.New("Please check the requestURL"), "Invalid requestURL")
	}

	//////////////////////////////////////// Proxy ////////////////////////////////////////////////////////////////
	// Notice: You should use proxy mode if you are in China, or you can not reach out Etherscan because of China Great Firewall.
	// If you don't need a proxy, you can delete it.
	// Note that I am using the default proxy port for Clash for Windows here: 127.0.0.1:7890
	proxyURL, _ := url.Parse("http://127.0.0.1:7890")

	// Create an HTTP client with a proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
	}

	// NOTICE: If you don't need a proxy client, you should use this code:
	// response, err := http.Get(requestURL)

	// Wait and retry if fail to get data from Etherscan
	var response *http.Response
	maxRetries := 5 // maximum number of retries
	for i := 0; i < maxRetries; i++ {
		// Sending GET requests using clients with proxies
		response, err = client.Get(requestURL)
		if err == nil {
			break // Success, exit loop
		}
		time.Sleep(1 * time.Second) // Wait for 1 second before retrying
	}
	defer response.Body.Close()
	if err != nil {
		log.Error("Timeout: Fail to fetch ABI from Etherscan. ChainID:", chainID, "contractAddress:", contractAddress)
		return nil, errors.Wrap(errors.New("Fail to fetch ABI from Etherscan"), "Timeout")
	}
	//////////////////////////////////////// Proxy ////////////////////////////////////////////////////////////////

	// Read response content
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Fail to Read response content. ChainID:", chainID, "contractAddress:", contractAddress)
		return []byte{}, errors.Wrap(errors.New("Fail to Read response content"), "Read response fail")
	}

	// Parsing JSON data
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Error("Fail to parsing JSON data. ChainID:", chainID, "contractAddress:", contractAddress)
		return []byte{}, errors.Wrap(errors.New("Fail to parsing JSON data"), "Parse fail")
	}

	return []byte(apiResponse.Result), nil
}

// @dev Query a contract's RuntimeCode
func queryRuntimeCode(rpcUrl string, contractAddress common.Address) ([]byte, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Error("Fail to connect to the node. RPC URL:", rpcUrl, "ContractAddress:", contractAddress)
		return nil, errors.Wrap(errors.New("Fail to connect to the node"), "Connect fail")
	}

	bytecode, err := client.CodeAt(context.Background(), contractAddress, nil) // nil: the newest block
	if err != nil {
		log.Error("Fail to get the RuntimeCode. RPC URL:", rpcUrl, "ContractAddress:", contractAddress)
		return nil, errors.Wrap(errors.New("Fail to get the RuntimeCode"), "Get fail")
	}

	if len(bytecode) == 0 {
		return []byte{}, nil
	} else {
		return bytecode, nil
	}

}
