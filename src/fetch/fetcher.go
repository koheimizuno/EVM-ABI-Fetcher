package fetch

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
)

// @dev For parse the data from Etherscan
type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type FetcherCli struct {
	ApiKey string
}

func (f *FetcherCli) GetABIAtStartOfBlock(db *gorm.DB, chainID int, contractAddress common.Address, block *big.Int) ([]byte, error) {

	var urlString string
	if chainID == 1 { // Ethereum
		urlString = fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", contractAddress, f.ApiKey)
	} else if chainID == 56 { // BSC

	} else if chainID == 42161 { // Arbitrum

	} else if chainID == 8453 { // Base

	} else if chainID == 43114 { // Avalanche

	} else if chainID == 137 { // Polygon

	} else {
		fmt.Println("Invalid chainID")
	}

	// Notice: You should use proxy mode if you are in China, or you can not reach out Etherscan because of China Great Firewall.
	// If you don't need a proxy, you can delete it.
	// Note that I am using the default proxy port for Clash for Windows here: 127.0.0.1:7890
	proxyURL, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		fmt.Printf("Fail to parse the proxy address: %s\n", err)
		return []byte("TODO"), nil
	}

	// Create an HTTP client with a proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
	}

	// Sending GET requests using clients with proxies
	response, err := client.Get(urlString)
	if err != nil {
		fmt.Printf("request fail: %s\n", err)
		return []byte("TODO"), nil
	}
	defer response.Body.Close()

	// Read response content
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("read the response fail: %s\n", err)
		return []byte("TODO"), nil
	}

	// Parsing JSON data
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Printf("Parsing JSON data fail: %s\n", err)
		return []byte("TODO"), nil
	}

	// Print the content of the result field
	fmt.Println("Result from Etherscan API:")
	fmt.Println(apiResponse.Result)

	// 1. Check if the ABI exists in the database for the given chainID, contract address, and function signature.
	//    If found, return the ABI from the database

	// 2. If the ABi is not found in the database, query Etherscan to retrieve the ABI
	//    If Etherscan returns the ABI, store it in the database and return it

	// 3. If Etherscan does not have the ABI, return an appropriate error

	return []byte("TODO"), nil
}
