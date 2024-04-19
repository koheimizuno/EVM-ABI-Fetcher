package fetch

import (
	myCache "code/src/cache"
	myDB "code/src/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/petermattis/goid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	ApiKey string // Etherscan
	RpcUrl string // Blockchain node RPC
	mu     sync.RWMutex
}

var cache = *myCache.NewABICache()
var log = logrus.New()
var f = FetcherCli{ApiKey: os.Getenv("API_KEY"), RpcUrl: os.Getenv("RPC_URL")}

var db = myDB.InitDatabase()

// GetFunctionABIAtBlock
// @dev try to get the function ABI
func GetFunctionABIAtBlock(chainID int, contractAddress common.Address, sig [4]byte, block *big.Int) (*abi.Method, error) {
	// [1. In memory]
	functionABI, _, isFound := cache.Get(chainID, contractAddress, string(sig[:]))
	if isFound {
		log.Info("[Thread ", goid.Get(), "] Found functionABI in cache, data:", functionABI)
		return functionABI, nil
	}

	// [2. In DB] Check if the functionABI exists in the database for the given chainID, contract address and sig
	ID := myCache.CacheKey(chainID, contractAddress, string(sig[:]))
	var functionSignature myDB.FunctionSignature
	if err := db.Where("id = ?", ID).First(&functionSignature).Error; err != nil { // Not found ABI in DB
		log.Error("Not found the functionABI in DB")

		f.mu.Lock()
		defer f.mu.Unlock()
		// logic: Not found the ABI in DB => if there is a shouldEtherscan item in DB?
		//           1. no: create a new shouldEtherscan item for the given chainID and contractAddress
		//           2. yes: check that whether now passes 2 days since the last time or not?
		//                1. no: do nothing
		//                2. yes: set searchEtherscan to true, then other thread of searchInEtherscan() will search from Etherscan

		var searchEtherscan myDB.SearchEtherscan
		now := time.Now().Unix()

		// search in DB
		result := db.Where("chain_id = ? AND contract_address = ?", chainID, contractAddress).First(&searchEtherscan)
		if result.Error != nil { // not found the searchEtherscan item by chainID nad contractAddress in DB
			// create a new item
			newRecord := myDB.SearchEtherscan{
				ChainID:         chainID,
				ContractAddress: contractAddress.Bytes(),
				Time:            int(now),
				ShouldSearch:    true, // should search in Etherscan
			}
			err = db.Create(&newRecord).Error
			if err != nil {
				log.Error("Fail to create a searchEtherscan item in db")
				return nil, errors.Wrap(errors.New("Fail to create an item in db"), "Create fail")
			}
		} else { // the record exists
			if now-int64(searchEtherscan.Time) >= 48*time.Hour.Microseconds() { // has pass 2 days?
				// pass 2 days, update shouldSearch to true. so the robot will search ABi from Etherscan by searchInEtherscan()
				err = db.Model(&searchEtherscan).Update("should_search", true).Error
				if err != nil {
					log.Error("Fail to update the searchEtherscan item to true in db")
					return nil, errors.Wrap(errors.New("Fail to update the item in db"), "Update fail")
				}
			}
		}
		log.Warning("Waiting robot to search the functionABI from Etherscan")
		return nil, errors.Wrap(errors.New("Waiting robot to search the functionABI from Etherscan"), "Not Found")
	} else { // found in db
		log.Info("Found functionABI in DB")

		///////////////////////////// update the cache /////////////////////////////////////////
		f.mu.Lock()
		defer f.mu.Unlock()

		// Second check
		functionABISecondCheck, _, isFoundSecondCheck := cache.Get(chainID, contractAddress, string(sig[:]))
		if isFoundSecondCheck { // If found functionABI in cache
			log.Info("[Thread ", goid.Get(), "] Second check found functionABI in cache")
			return functionABISecondCheck, nil
		} else { // not found in cache, set the cache

			// define the data to search in DB
			var resultContractABIID = functionSignature.ContractBytecodeID
			var contractBytecode myDB.ContractBytecode
			_ = db.Where("id = ?", resultContractABIID).First(&contractBytecode)
			log.Info("hello")
			// unmarshal the functionABI
			myABI, err := abi.JSON(strings.NewReader(functionSignature.FunctionABI))
			if err != nil {
				log.Info("Fail to unmarshal the ABI")
				return nil, errors.Wrap(errors.New("Fail to unmarshal the ABI"), "Unmarshal fail")
			}
			// get the Method's key, then we can use the key to find the functionABI(type: abi.Method)
			//////////////////////////////////////////////////////////////////////
			var resultFunctonABI abi.Method
			for key := range myABI.Methods {
				// get the functionABI(type: abi.Method) by key
				resultFunctonABI = myABI.Methods[key] // ensure the variable to be marshaled is abi.Method
			}

			var resultContractABI *abi.ABI

			err = json.Unmarshal([]byte(contractBytecode.ContractABI), &resultContractABI)
			if err != nil {
				log.Error("Fail to unmarshal ContractABI. Err:", err)
				return nil, errors.Wrap(errors.New("Fail to unmarshal ContractABI"), "Fail unmarshal")
			}

			// set the data to cache
			cache.Set(
				chainID,
				contractAddress,
				&resultFunctonABI,
				resultContractABI,
				string(functionSignature.Signature),
			)
			///////////////////////////// update the cache /////////////////////////////////////////

			return &resultFunctonABI, nil // return the functionABI from DB
		}

	}

}

// @dev Set up some robot threads to run this function, search ABI from Etherscan
func searchInEtherscan(apiKey string, rpcUrl string) error {

	// 1.Update the shouldSearch field
	var resultsFalse []myDB.SearchEtherscan
	// If the item pass 2 days, we set the shouldSearch field to true, so that it will try to get the ABi from Etherscan
	err := db.Where("should_search = ?", false).Find(&resultsFalse).Error
	if err != nil {
		log.Error("Fail to search item in DB")
		return errors.Wrap(errors.New("Fail to search item in DB"), "Search fail")
	}
	for _, item := range resultsFalse {
		// has pass 2 days => update the shouldSearch to true
		if time.Now().Unix()-int64(item.Time) >= 48*time.Hour.Microseconds() {
			err = db.Model(&item).Update("should_search", true).Error
			if err != nil {
				log.Error("Fail to update the searchEtherscan item in db")
				return errors.Wrap(errors.New("Fail to update the item in db"), "Update fail")
			}
		}
	}

	// 2.Iterator the DB, if the shouldSearch field is true, than search ABI in Etherscan
	f.mu.Lock()
	defer f.mu.Unlock()

	var results []myDB.SearchEtherscan
	// query the records: shouldSearch = true
	err = db.Where("should_search = ?", true).Find(&results).Error
	if err != nil {
		log.Error("Fail to search item in db")
		return errors.Wrap(errors.New("Fail to search item in db"), "Search fail")
	}

	// 3.The all items whose shouldSearch field are true
	for _, item := range results {
		log.Info("Begin search ABI from Etherscan. ChinaID:", item.ChainID, ".contractAddress:", item.ContractAddress)

		var contractAddress common.Address
		copy(contractAddress[:], item.ContractAddress[:])

		// Begin search ABI in Etherscan
		data, err := queryABIFromEtherscan(apiKey, item.ChainID, contractAddress)
		if err != nil {
			log.Error("Fail to search item in Etherscan")
			return errors.Wrap(errors.New("Fail to search item in Etherscan"), "Search fail")
		}

		// Begin search Bytecode in blockchain node
		bytecode, err := queryRuntimeCode(rpcUrl, contractAddress)
		if err != nil {
			log.Error("Fail to search bytecode")
			return errors.Wrap(errors.New("Fail to search bytecode"), "Search fail")
		}
		// store the contract's info into DB. [bytecode]
		contractbytecodId := uuid.New()
		contract := myDB.ContractBytecode{
			ID:                contractbytecodId,
			Bytecode:          bytecode,     // the contract's bytecode
			SourceCode:        "",           // TODO
			CompileTimeParams: "",           // TODO
			ContractABI:       string(data), // the contract's ABI
		}
		err = db.Create(&contract).Error
		if err != nil {
			log.Error("Fail to create an item")
			return errors.Wrap(errors.New("Fail to create an item"), "Create fail")
		}

		// Using JSON RawMessage to maintain the original JSON format
		var rawMessages []json.RawMessage
		_ = json.Unmarshal(data, &rawMessages)
		// Create a new string array to store each object
		var functionStrings []string
		for _, raw := range rawMessages {
			functionStrings = append(functionStrings, string(raw))
		}

		for _, funcStr := range functionStrings {
			theABI, err := abi.JSON(strings.NewReader("[" + funcStr + "]"))
			if err != nil {
				log.Error("Fail to parse the abi")
				return errors.Wrap(errors.New("Fail to parse the abi"), "Parse fail")
			} else {

				// get the Method's key, then we can use the key to find the functionABI(type: abi.Method)
				for key := range theABI.Methods {
					// get the functionABI(type: abi.Method) by key
					function := theABI.Methods[key] // ensure the variable to be marshaled is abi.Method
					// functionABI(type: abi.Method) => signature => 4bytes signature
					sig4bytes := myCache.Get4bytesSig(function.Sig)

					// set the functionABI to DB
					ID := myCache.CacheKey(item.ChainID, contractAddress, string(sig4bytes[:]))
					// TODO: marshal the functionABI, later it fails to unmarshal
					functionSig := myDB.FunctionSignature{
						ID:                 ID,
						ContractBytecodeID: contractbytecodId,
						Signature:          sig4bytes[:],
						FunctionABI:        "[" + funcStr + "]",
					}
					err = db.Create(&functionSig).Error
					if err != nil {
						log.Error("Fail to create a FunctionSignature item")
						return errors.Wrap(errors.New("Fail to create a FunctionSignature item"), "Create fail")
					}
				}
				// After get the ABI from Etherscan, set the shouldSearch to false
				result := db.Model(&myDB.SearchEtherscan{}).
					Where("chain_id = ? AND contract_address = ?", item.ChainID, item.ContractAddress).
					Update("should_search", false)
				if result.Error != nil {
					log.Error("Fail to create a FunctionSignature item")
					return errors.Wrap(errors.New("Fail to update the shouldSearch field"), "Update fail")
				}
			}
		}
	}

	return nil
}

// GetABIAtStartOfBlock
// @notice We do not use the block parameter for now
// @param db The SQLite3's handle
// @param chainID The chain ID
// @param contractAddress The contract whose ABI you want to get
// @param block Get the ABI in a certain block height.
// TODO old version
/*
func (f *FetcherCli) GetContractABIAtBlock(db *gorm.DB, chainID int, contractAddress common.Address, block *big.Int) ([]byte, error) {

	// [1. In memory]
	_, contractABI, isFound := cache.Get(chainID, contractAddress, "Search for contractABI")
	if isFound {
		log.Info("Thread ", goid.Get(), ". Found contractABI in cache, length:", len(contractABI))
		return contractABI, nil
	}

	// [2. In DB] Check if the contractABI exists in the database for the given chainID and contract address
	//               If found, return the contractABI from the database

	addressBytes := contractAddress.Bytes() // Convert Ethereum addresses to a format that can be queried by the database

	var deployment myDB.ContractDeployment
	// Query in the database
	result := db.Where("chain_id = ? AND contract_address = ?", chainID, addressBytes).First(&deployment)
	if result.Error == nil {
		// Search for the specific contractABI by ContractDeployment.ContractBytecodeID
		var bytecode myDB.ContractBytecode
		result = db.Where("id = ?", deployment.ContractBytecodeID).First(&bytecode)
		if result.Error == nil { // result.Error is equal to nil means that there is an contractABI in the database
			///////////////////// Write Lock //////////////////////////
			f.mu.Lock()
			defer f.mu.Unlock()

			// Second check to prevent this situation: Multiple threads complete read operations simultaneously and then queue up for write operations,
			// which can cause duplicate cache writes
			_, contractABISecondCheck, isFoundSecondCheck := cache.Get(chainID, contractAddress, "Search for contractABI")
			if isFoundSecondCheck {
				log.Info("Thread ", goid.Get(), ". Found contractABISecondCheck in cache, length:", len(contractABISecondCheck))

				return contractABISecondCheck, nil
			} else {
				// Set to cache: We search for contractABI, so this cache item will not contain functionABI or signature
				cache.Set(
					chainID,
					contractAddress,
					[]byte("Nil because setting the cache when search for contractABI in DB"),
					[]byte(bytecode.ContractABI),
					"Nil because setting the cache when search for contractABI in DB",
				)
			}

			///////////////////// Write Lock //////////////////////////

			log.Info("Thread ", goid.Get(), ". Found ContractABI in DB, length:", len(bytecode.ContractABI))
			return []byte(bytecode.ContractABI), nil
		}
	}

	// [3. Etherscan] If the ABi is not found in the database, query Etherscan to retrieve the ABI
	//              If Etherscan returns the ABI, store it in the database and return it

	// Create a new record of ContractBytecode
	bytecode, err := queryRuntimeCode(f.RpcUrl, contractAddress)
	if string(bytecode) == "" {
		log.Error("Thread ", goid.Get(), ".The address doesn't contain RuntimeCode")
		return nil, errors.Wrap(errors.New("Check the address you input"), "Not contract")
	}

	// If tht contract has not been verified, it will return: Contract source code not verified
	abi, err := queryABIFromEtherscan(f.ApiKey, chainID, contractAddress)
	if string(abi) == "Contract source code not verified" { // EOA or the not verified contract
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
		ID:                contractBytecodeID,
		Bytecode:          bytecode,
		SourceCode:        "TODO", // TODO
		CompileTimeParams: "TODO", // TODO
		ContractABI:       "TODO", // TODO
	}
	db.Create(&newContractBytecode)

	// TODO 使用for循环，将所有的函数签名都分别插入到这个表当中
	// Store the ABI that fetched from Etherscan into the database
	newSignature := myDB.FunctionSignature{
		ID:                 functionSignatureID,
		ContractBytecodeID: contractBytecodeID,
		Signature:          nil,    // TODO How to deal with this， what to store: single signature or a array of signatures
		FunctionABI:        "TODO", // TODO
	}
	db.Create(&newSignature)

	// Create a new record of ContractDeployment
	newDeployment := myDB.ContractDeployment{
		ChainID:            chainID,
		ContractAddress:    addressBytes,
		ContractBytecodeID: contractBytecodeID,
	}
	db.Create(&newDeployment)

	// write to cache

	_, contractABI, isFound = cache.Get(chainID, contractAddress, "Search for contractABI")
	if isFound {
		log.Info("Thread ", goid.Get(), ". Found ABI in cache, length:", len(contractABI))
		return contractABI, nil
	} else {
		cache.Set(
			chainID,
			contractAddress,
			[]byte("Nil because setting the cache when search for contractABI in DB"),
			contractABI,
			"Nil because setting the cache when search for contractABI in DB",
		)
	}

	///////////////////// Write Lock //////////////////////////

	return abi, nil
}
*/

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
		requestURL = fmt.Sprintf("https://api.bscscan.com/api?module=contract&action=getabi&address=%s&apikey=%s", contractAddress, apiKey)
	} else if chainID == 42161 { // Arbitrum
		requestURL = fmt.Sprintf("https://api.arbiscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", contractAddress, apiKey)
	} else if chainID == 137 { // Polygon
		requestURL = fmt.Sprintf("https://api.polygonscan.com/api?module=contract&action=getabi&address=%s&apikey=%s", contractAddress, apiKey)
	} else {
		log.Error("Invalid chainID or API KEY. chainID:", chainID, "API KEY:", apiKey, "contractAddress", contractAddress)
		return "", errors.Wrap(errors.New("Check the input"), "Fail to checkChainIDAndGetReqURL")
	}

	return requestURL, nil
}

// @dev Query a contract's ABI from Etherscan
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

// @dev Query a contract's RuntimeCode from block node
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
