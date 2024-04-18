package fetch

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var contractAddress1 = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") // USDT
var contractAddress2 = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") // WETH
var contractAddress3 = common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F") // DAI
var signature1 = [4]byte{142, 88, 167, 150}
var signature2 = [4]byte{220, 163, 100, 249}
var signature3 = [4]byte{193, 190, 199, 214}
var blockHeight = big.NewInt(10000)

func TestGetABIAtStartOfBlockInOneThread_ContainEOA(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	//apiKey := os.Getenv("API_KEY")
	//rpcURL := os.Getenv("RPC_URL")

	data, err := GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	fmt.Println("data:", data)
	_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	_, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	assert.NoError(t, err)
	//
	//// 创建一个ticker，每21秒触发一次
	//ticker := time.NewTicker(5 * time.Second)
	//// 创建一个计时器，25秒后触发
	//stopTimer := time.NewTimer(8 * time.Second)
	//
	//for {
	//	select {
	//	case <-ticker.C: // 每当ticker发出信号时
	//		_ = searchInEtherscan(apiKey, rpcURL)
	//	case <-stopTimer.C: // 当计时器结束时
	//		ticker.Stop() // 停止ticker
	//
	//		/////////////////////////////// Found in db //////////////////////////////////
	//		_, err = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	//		_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	//		_, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	//		/////////////////////////////// Found in db //////////////////////////////////
	//
	//		/////////////////////////////// Found in cache //////////////////////////////////
	//		_, err = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	//		_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	//		_, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	//
	//		_, err = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	//		_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	//		_, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	//
	//		_, err = GetFunctionABIAtBlock(1, contractAddress1, signature1, blockHeight)
	//		_, err = GetFunctionABIAtBlock(1, contractAddress2, signature2, blockHeight)
	//		data, err = GetFunctionABIAtBlock(1, contractAddress3, signature3, blockHeight)
	//		fmt.Println("data", data)
	//		/////////////////////////////// Found in cache //////////////////////////////////
	//		fmt.Println("stop")
	//		return // 退出循环和程序
	//	}
	//}

}

/*

// Write unit test using the `testing` package in Go to cover the core functionality of the GetABI function

// 2.42s
func TestGetABIAtStartOfBlockInOneThread_ContainEOA(t *testing.T) {

	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")
	rpcURL := os.Getenv("RPC_URL")

	db := myDB.InitDatabase()

	fetcher := FetcherCli{ApiKey: apiKey, RpcUrl: rpcURL}
	fmt.Println("First search")
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")) // Verify
	assert.NoError(t, err)

	fmt.Println("Second search")                                                                                    // Test the caching mechanism to ensure that frequently accessed ABIs are retrieved from the cache
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")) // Verify
	assert.NoError(t, err)

	fmt.Println("Third search")                                                                                     // Test the caching mechanism to ensure that frequently accessed ABIs are retrieved from the cache
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")) // Verify
	assert.NoError(t, err)

	fmt.Println("Fourth search")                                                                                    // Test the caching mechanism to ensure that frequently accessed ABIs are retrieved from the cache
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")) // Verify
	assert.NoError(t, err)

	fmt.Println("Fifth search")                                                                                     // Test the caching mechanism to ensure that frequently accessed ABIs are retrieved from the cache
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")) // Verify
	assert.NoError(t, err)

	fmt.Println("Six search")
	data, err := fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xd3E65149C212902749D49011B6ab24bba30D97c6")) // EOA
	assert.Equal(t, []byte(nil), data)
	assert.Error(t, err)

	fmt.Println("Seven search")
	data, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xf88d1d6d9db9a39dbbfc4b101cecc495bb0636f8")) // Not verify
	assert.Equal(t, []byte(nil), data)
	assert.Error(t, err)

	var deployment myDB.ContractDeployment
	result := db.Where("chain_id = ? AND contract_address = ?", 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")).First(&deployment)
	if result.Error == nil {
		var signature myDB.FunctionSignature
		result = db.Where("id = ?", deployment.FunctionSignatureID).First(&signature)
		if result.Error == nil { // Check that the ABI of 0xdAC17F958D2ee523a2206206994597C13D831ec7(USDT) is correct
			assert.Equal(t, signature.FunctionABI, "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_upgradedAddress\",\"type\":\"address\"}],\"name\":\"deprecate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"deprecated\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_evilUser\",\"type\":\"address\"}],\"name\":\"addBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradedAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maximumFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_maker\",\"type\":\"address\"}],\"name\":\"getBlackListStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBasisPoints\",\"type\":\"uint256\"},{\"name\":\"newMaxFee\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"basisPointsRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isBlackListed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_clearedUser\",\"type\":\"address\"}],\"name\":\"removeBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_UINT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_blackListedUser\",\"type\":\"address\"}],\"name\":\"destroyBlackFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_initialSupply\",\"type\":\"uint256\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_decimals\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Issue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Redeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"Deprecate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"feeBasisPoints\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"maxFee\",\"type\":\"uint256\"}],\"name\":\"Params\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_blackListedUser\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"DestroyedBlackFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"AddedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"RemovedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]")
		}
		if result.Error == nil { // Check that the ABI of 0xdAC17F958D2ee523a2206206994597C13D831ec7(USDT) is correct
			assert.NotEqual(t, signature.FunctionABI, "test [Handle any discrepancies or inconsistencies between the retrieved ABIs and the expected ABIs]")
		}
	}
}

// 2.09s
// Write integration test to verify the interaction between the GetABI function, the database, and Etherscan
func TestGetABIAtStartOfBlockInMultipleThreads_ContainEOA(t *testing.T) {

	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")
	rpcURL := os.Getenv("RPC_URL")

	db := myDB.InitDatabase()

	addressArray1 := [5]string{
		"0xd3E65149C212902749D49011B6ab24bba30D97c6", // EOA 1
		"0x63D248D9F6562a7b7a76ca485CC3564aFA9cbd00", // EOA 2
		"0x6a8C95fFAdAf369A73F423A464d334Aa6158e259", // EOA 3
		"0x4536c8d46eC3C127A0868ea9dE83D7dd5e01e0E3", // EOA 4
		"0xa29E4fe451CCFa5e7DEF35188919ad7077A4DE8f", // contract 2: not verify
	}

	var wg1 sync.WaitGroup
	wg1.Add(len(addressArray1))
	fetcher1 := FetcherCli{ApiKey: apiKey, RpcUrl: rpcURL}
	for _, address := range addressArray1 {
		go func(addr string) {
			defer wg1.Done()
			_, err = fetcher1.GetABIAtStartOfBlock(db, 1, common.HexToAddress(addr))
			// Test error handling for scenarios where the ABI is not found or Etherscan returns an error
			assert.Error(t, err)
		}(address)
	}

	wg1.Wait()

	addressArray2 := [5]string{
		"0xd9e1cE17f2641f24aE83637ab66a2cca9C378B9F", // contract 3: verify, SushiSwap V2 router
		"0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f", // contract 4: verify, Uniswap V2 factory
		"0xdAC17F958D2ee523a2206206994597C13D831ec7", // contract 5: verify, USDT
		"0x6B175474E89094C44Da98b954EedeAC495271d0F", // contract 6: verify, DAI
		"0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", // contract 8: verify, Uniswap V2 router
	}

	var wg2 sync.WaitGroup
	wg2.Add(len(addressArray2))
	fetcher := FetcherCli{ApiKey: apiKey, RpcUrl: rpcURL}
	for _, address := range addressArray2 {
		go func(addr string) {
			defer wg2.Done()
			_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress(addr))
			assert.NoError(t, err)
		}(address)
	}

	wg2.Wait()

	// Test scenarios where the ABI is found in the database and retrieved from Etherscan
	var deployment myDB.ContractDeployment
	result := db.Where("chain_id = ? AND contract_address = ?", 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")).First(&deployment)
	if result.Error == nil {
		var signature myDB.FunctionSignature
		result = db.Where("id = ?", deployment.FunctionSignatureID).First(&signature)
		if result.Error == nil { // Compare the retrieved ABIs with the expected ABIs from reliable sources to ensure correctness
			assert.Equal(t, signature.FunctionABI, "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_upgradedAddress\",\"type\":\"address\"}],\"name\":\"deprecate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"deprecated\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_evilUser\",\"type\":\"address\"}],\"name\":\"addBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradedAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maximumFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_maker\",\"type\":\"address\"}],\"name\":\"getBlackListStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBasisPoints\",\"type\":\"uint256\"},{\"name\":\"newMaxFee\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"basisPointsRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isBlackListed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_clearedUser\",\"type\":\"address\"}],\"name\":\"removeBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_UINT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_blackListedUser\",\"type\":\"address\"}],\"name\":\"destroyBlackFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_initialSupply\",\"type\":\"uint256\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_decimals\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Issue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Redeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"Deprecate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"feeBasisPoints\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"maxFee\",\"type\":\"uint256\"}],\"name\":\"Params\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_blackListedUser\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"DestroyedBlackFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"AddedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"RemovedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]")
		}
	}

	// There should be 5 items that have bytecode and ABI

	var count int64
	db.Model(&myDB.ContractDeployment{}).Count(&count)
	assert.Equal(t, count, int64(5))

}

// 0.03s Starting from the second test, the testing time is 0.03s because the data has already been stored in the database
func TestGetABIAtStartOfBlockInMultipleThreads_NotContainEOA(t *testing.T) {

	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")
	rpcURL := os.Getenv("RPC_URL")

	db := myDB.InitDatabase()

	addressArray1 := [2]string{
		"0xdAC17F958D2ee523a2206206994597C13D831ec7", // contract 1: verify, USDT
		"0xd9e1cE17f2641f24aE83637ab66a2cca9C378B9F", // contract 2: verify, SushiSwap V2 router
	}

	var wg1 sync.WaitGroup
	wg1.Add(len(addressArray1))
	fetcher1 := FetcherCli{ApiKey: apiKey, RpcUrl: rpcURL}
	for _, address := range addressArray1 {
		go func(addr string) {
			defer wg1.Done()
			_, err = fetcher1.GetABIAtStartOfBlock(db, 1, common.HexToAddress(addr))
			assert.NoError(t, err)
		}(address)
	}

	wg1.Wait()

	addressArray2 := [4]string{
		"0x6B175474E89094C44Da98b954EedeAC495271d0F", // contract 3: verify, DAI
		"0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f", // contract 4: verify, Uniswap V2 factory
		"0xdAC17F958D2ee523a2206206994597C13D831ec7", // contract 7: verify, USDT
		"0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", // contract 8: verify, Uniswap V2 router
	}

	var wg2 sync.WaitGroup
	wg2.Add(len(addressArray2))
	fetcher := FetcherCli{ApiKey: apiKey, RpcUrl: rpcURL}
	for _, address := range addressArray2 {
		go func(addr string) {
			defer wg2.Done()
			_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress(addr))
			assert.NoError(t, err)
		}(address)
	}

	wg2.Wait()

	var deployment myDB.ContractDeployment
	result := db.Where("chain_id = ? AND contract_address = ?", 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")).First(&deployment)
	if result.Error == nil {
		var signature myDB.FunctionSignature
		result = db.Where("id = ?", deployment.FunctionSignatureID).First(&signature)
		if result.Error == nil { // Check that the ABI of 0xdAC17F958D2ee523a2206206994597C13D831ec7(USDT) is correct
			assert.Equal(t, signature.FunctionABI, "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_upgradedAddress\",\"type\":\"address\"}],\"name\":\"deprecate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"deprecated\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_evilUser\",\"type\":\"address\"}],\"name\":\"addBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradedAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maximumFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_maker\",\"type\":\"address\"}],\"name\":\"getBlackListStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBasisPoints\",\"type\":\"uint256\"},{\"name\":\"newMaxFee\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"basisPointsRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isBlackListed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_clearedUser\",\"type\":\"address\"}],\"name\":\"removeBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_UINT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_blackListedUser\",\"type\":\"address\"}],\"name\":\"destroyBlackFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_initialSupply\",\"type\":\"uint256\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_decimals\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Issue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Redeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"Deprecate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"feeBasisPoints\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"maxFee\",\"type\":\"uint256\"}],\"name\":\"Params\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_blackListedUser\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"DestroyedBlackFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"AddedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"RemovedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]")
		}
	}

	// There should be 5 items that have bytecode and ABI
	var count int64
	db.Model(&myDB.ContractDeployment{}).Count(&count)
	assert.Equal(t, count, int64(5))
}

func TestCheckChainIDAndGetReqURL(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")
	correctURL := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=0xdAC17F958D2ee523a2206206994597C13D831ec7&apikey=%s", apiKey)

	url, err := checkChainIDAndGetReqURL(apiKey, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"))
	assert.Equal(t, correctURL, url)

	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	// Valid Ethereum chain ID
	t.Run("Valid Ethereum Chain ID", func(t *testing.T) {
		url, err := checkChainIDAndGetReqURL(apiKey, 1, contractAddress)
		require.NoError(t, err)
		assert.Contains(t, url, "https://api.etherscan.io/api")
	})

	// Invalid Chain ID
	t.Run("Invalid Chain ID", func(t *testing.T) {
		_, err := checkChainIDAndGetReqURL(apiKey, 9999, contractAddress)
		assert.Error(t, err)
	})

}

func TestQueryABIFromEtherscan(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	apiKey := os.Getenv("API_KEY")

	data, err := queryABIFromEtherscan(apiKey, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")) // Contract
	assert.NoError(t, err)
	assert.Greater(t, len(data), 0)

	data, err = queryABIFromEtherscan(apiKey, 1, common.HexToAddress("0xd3E65149C212902749D49011B6ab24bba30D97c6")) // EOA
	assert.NoError(t, err)
	assert.Equal(t, string(data), "Contract source code not verified")
}

func TestQueryRuntimeCode(t *testing.T) {
	err := godotenv.Load("../../.env") // Load the `.env` file in the current directory by default
	assert.NoError(t, err)

	rpcURL := os.Getenv("RPC_URL")

	data, err := queryRuntimeCode(rpcURL, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")) // Contract
	assert.NoError(t, err)
	assert.Greater(t, len(data), 0)

	data, err = queryRuntimeCode(rpcURL, common.HexToAddress("0xd3E65149C212902749D49011B6ab24bba30D97c6")) // EOA
	assert.NoError(t, err)
	assert.Equal(t, 0, len(data))
}

*/
