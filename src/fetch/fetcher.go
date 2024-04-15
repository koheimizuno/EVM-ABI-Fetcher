package fetch

import (
	myDB "code/src/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
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
}

// GetABIAtStartOfBlock
// @notice We do not use the block parameter for now
// @param db The SQLite3's handle
// @param chainID The chain ID
// @param contractAddress The contract whose ABI you want to get
// @param block Get the ABI in a certain block height.
//
//	NOTICE: There is no such an API in the Etherscan, so this parameter won't be used temporarily
func (f *FetcherCli) GetABIAtStartOfBlock(db *gorm.DB, chainID int, contractAddress common.Address) ([]byte, error) {

	// 1. Check if the ABI exists in the database for the given chainID, contract address, and function signature. If found, return the ABI from the database

	addressBytes := contractAddress.Bytes() // Convert Ethereum addresses to a format that can be queried by the database

	var deployment myDB.ContractDeployment
	// Query in the database
	result := db.Where("chain_id = ? AND contract_address = ?", chainID, addressBytes).First(&deployment)
	if result.Error == nil {
		// Search for the specific ABI by ContractDeployment.FunctionSignatureID
		var signature myDB.FunctionSignature
		result := db.Where("id = ?", deployment.FunctionSignatureID).First(&signature)
		if result.Error == nil { // result.Error is equal to nil means that there is an ABI in the database of the input contractAddress
			return []byte(signature.FunctionABI), nil
		}
	}

	// 2. If the ABi is not found in the database, query Etherscan to retrieve the ABI
	//    If Etherscan returns the ABI, store it in the database and return it

	abi, err := queryABIFromEtherscan(f.ApiKey, chainID, contractAddress)
	// 3. If Etherscan does not have the ABI, return an appropriate error
	if err != nil {
		return nil, errors.Wrap(errors.New("Fail to fetch ABI by Etherscan API"), "Fetch fail")
	}

	// Prepare some UUID
	functionSignatureID := uuid.New()
	contractBytecodeID := uuid.New()

	// Create a new record of ContractBytecode
	bytecode, err := queryRuntimeCode(f.RpcUrl, chainID, contractAddress)
	if err != nil {
		return nil, errors.Wrap(errors.New("Fail to fetch Bytecode by Etherscan API"), "Fetch fail")
	}

	newContractBytecode := myDB.ContractBytecode{
		ID:       contractBytecodeID,
		Bytecode: bytecode,
	}
	db.Create(&newContractBytecode)

	// Store the ABI that fetched from Etherscan into the database
	newSignature := myDB.FunctionSignature{
		ID:          functionSignatureID,
		Signature:   nil, // TODO How to deal with this
		FunctionABI: abi, // TODO How to deal with this
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

	return []byte(abi), nil
}

// @dev Check ChainID and get the format the request url
func checkChainIDAndGetReqURL(apiKey string, chainID int, contractAddress common.Address) (string, error) {
	var requestURL string

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
		return "", errors.Wrap(errors.New("You should provide a valid chainID"), "Invalid chainID")
	}

	return requestURL, nil
}

// @dev Query a contract's ABI
func queryABIFromEtherscan(apiKey string, chainID int, contractAddress common.Address) (string, error) {
	requestURL, err := checkChainIDAndGetReqURL(apiKey, chainID, contractAddress)
	if err != nil {
		return "", errors.Wrap(errors.New("Please check the chainID"), "Invalid chainID")
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

	// Sending GET requests using clients with proxies
	response, err := client.Get(requestURL)
	if err != nil {
		return "", errors.Wrap(errors.New("Fail to send a request"), "Request fail")
	}
	defer response.Body.Close()
	//////////////////////////////////////// Proxy ////////////////////////////////////////////////////////////////

	// Read response content
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(errors.New("Fail to get the response"), "Response fail")
	}

	// Parsing JSON data
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", errors.Wrap(errors.New("Fail to parsing JSON data"), "Parse fail")
	}

	return apiResponse.Result, nil
}

// @dev Query a contract's RuntimeCode
func queryRuntimeCode(rpcUrl string, chainID int, contractAddress common.Address) ([]byte, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, errors.Wrap(errors.New("Fail to connect to the node"), "Connect fail")
	}

	bytecode, err := client.CodeAt(context.Background(), contractAddress, nil) // nil: the newest block
	if err != nil {
		return nil, errors.Wrap(errors.New("Fail to get the RuntimeCode"), "Get fail")
	}

	if len(bytecode) == 0 {
		return []byte{}, nil
	} else {
		return bytecode, nil
	}

}
