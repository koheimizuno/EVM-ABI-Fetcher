package fetch

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"
)

var contractAddress1 = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") // USDT
var contractAddress2 = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") // WETH
var contractAddress3 = common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F") // DAI
var signature1 = [4]byte{142, 88, 167, 150}
var signature2 = [4]byte{220, 163, 100, 249}
var signature3 = [4]byte{193, 190, 199, 214}
var blockHeight = big.NewInt(10000)

func TestUnmarshal(t *testing.T) {
	jsonData := `
[
{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},
{"constant":false,"inputs":[{"name":"_upgradedAddress","type":"address"}],"name":"deprecate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},
{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},
{"constant":true,"inputs":[],"name":"deprecated","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"}]`
	var rawMessages []json.RawMessage
	err := json.Unmarshal([]byte(jsonData), &rawMessages)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	var functionStrings []string
	for _, raw := range rawMessages {
		functionStrings = append(functionStrings, string(raw))
	}

	for _, funcStr := range functionStrings {
		abi, err := abi.JSON(strings.NewReader("[" + funcStr + "]"))
		if err != nil {
			fmt.Println("err:", err)
		} else {
			fmt.Println("abi:", abi)
		}
	}
}

// Test get functionABI
func TestFetchFunctionABIFromEtherscan(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")
	rpcURL := os.Getenv("RPC_URL")
	_, _ = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	_ = searchInEtherscan(apiKey, rpcURL) // search ABI from Etherscan
	data, _ := GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	fmt.Println("the data[ name() ]:", data)
}

// Test get functionABI
func TestFetchContractABIFromEtherscan(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")
	rpcURL := os.Getenv("RPC_URL")
	_, _ = GetContractABIAtBlock(1, contractAddress1, blockHeight)
	_ = searchInEtherscan(apiKey, rpcURL) // search ABI from Etherscan
	data, _ := GetContractABIAtBlock(1, contractAddress1, blockHeight)
	fmt.Println("the data[ name() ]:", data)
}

// Test get contractABI
func TestGetFunctionABIAtBlock_InOneThread(t *testing.T) {

	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")
	rpcURL := os.Getenv("RPC_URL")

	// not found functionABI in cache and DB, then the shouldSearch field will be set to true.
	// then we can call searchInEtherscan() to search ABI from Etherscan
	data, err := GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	fmt.Println("data1:", data)
	data, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	fmt.Println("data2:", data)
	data, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	fmt.Println("data3:", data)

	_ = searchInEtherscan(apiKey, rpcURL) // search ABI from Etherscan
	fmt.Println()

	time1 := time.Now().Unix()
	/////////////////////////////// Found in db //////////////////////////////////
	_, err = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	fmt.Println()
	/////////////////////////////// Found in db //////////////////////////////////

	/////////////////////////////// Found in cache //////////////////////////////////
	_, err = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	fmt.Println()

	_, err = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	fmt.Println()

	_, err = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	/////////////////////////////// Found in cache //////////////////////////////////

	fmt.Println()
	time2 := time.Now().Unix()
	runtime := time1 - time2
	fmt.Println("search data in DB and cache, time:", runtime, "seconds")
	return
}

// Test get contractABI
func TestGetContractABIAtBlock_InOneThread(t *testing.T) {

	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")
	rpcURL := os.Getenv("RPC_URL")

	// not found functionABI in cache and DB, then the shouldSearch field will be set to true.
	// then we can call searchInEtherscan() to search ABI from Etherscan
	data, err := GetContractABIAtBlock(1, contractAddress1, blockHeight)
	fmt.Println("data1:", data)
	data, err = GetContractABIAtBlock(1, contractAddress2, blockHeight)
	fmt.Println("data2:", data)
	data, err = GetContractABIAtBlock(1, contractAddress3, blockHeight)
	fmt.Println("data3:", data)

	_ = searchInEtherscan(apiKey, rpcURL) // search ABI from Etherscan
	fmt.Println()

	time1 := time.Now().Unix()
	/////////////////////////////// Found in db //////////////////////////////////
	_, err = GetContractABIAtBlock(1, contractAddress1, blockHeight)
	_, err = GetContractABIAtBlock(1, contractAddress2, blockHeight)
	_, err = GetContractABIAtBlock(1, contractAddress3, blockHeight)
	fmt.Println()
	/////////////////////////////// Found in db //////////////////////////////////

	/////////////////////////////// Found in cache //////////////////////////////////
	_, err = GetContractABIAtBlock(1, contractAddress1, blockHeight)
	_, err = GetContractABIAtBlock(1, contractAddress2, blockHeight)
	_, err = GetContractABIAtBlock(1, contractAddress3, blockHeight)
	fmt.Println()

	_, err = GetContractABIAtBlock(1, contractAddress1, blockHeight)
	_, err = GetContractABIAtBlock(1, contractAddress2, blockHeight)
	_, err = GetContractABIAtBlock(1, contractAddress3, blockHeight)
	fmt.Println()

	_, err = GetContractABIAtBlock(1, contractAddress1, blockHeight)
	_, err = GetContractABIAtBlock(1, contractAddress2, blockHeight)
	_, err = GetContractABIAtBlock(1, contractAddress3, blockHeight)
	/////////////////////////////// Found in cache //////////////////////////////////

	fmt.Println()
	time2 := time.Now().Unix()
	runtime := time1 - time2
	fmt.Println("search data in DB and cache, time:", runtime, "seconds")
	return
}

// Test checkChainIDAndReqURL
func TestCheckChainIDAndReqURL(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")

	url, err := checkChainIDAndGetReqURL(apiKey, 1, contractAddress1)
	correctReq := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", contractAddress1, apiKey)
	assert.NoError(t, err)
	assert.Equal(t, correctReq, url)
}

// Test queryABIFromEtherscan
func TestQueryABIFromEtherscan(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")

	abi, err := queryABIFromEtherscan(apiKey, 1, contractAddress1)
	assert.NoError(t, err)
	assert.NotNil(t, abi)
}

// Test queryRuntimeCode
func TestQueryRuntimeCode(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	rpcURL := os.Getenv("RPC_URL")

	abi, err := queryRuntimeCode(rpcURL, contractAddress1)
	assert.NoError(t, err)
	assert.NotNil(t, abi)
}
