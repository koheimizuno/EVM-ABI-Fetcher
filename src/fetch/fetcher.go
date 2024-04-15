package fetch

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io"
	"math/big"
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
}

// GetABIAtStartOfBlock
// @notice We do not use the block parameter for now
// @param db The SQLite3's handle
// @param chainID The chain ID
// @param contractAddress The contract whose ABI you want to get
// @param block Get the ABI in a certain block height.
//
//	NOTICE: There is no such an API in the Etherscan, so this parameter won't be used temporarily
func (f *FetcherCli) GetABIAtStartOfBlock(db *gorm.DB, chainID int, contractAddress common.Address, block *big.Int) ([]byte, error) {

	requestURL, err := checkChainIDAndGetReqURL(f.ApiKey, chainID, contractAddress)
	if err != nil {
		return nil, errors.Wrap(errors.New("Please check the chainID"), "Invalid chainID")
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
		return nil, errors.Wrap(errors.New("Fail to send a request"), "Request fail")
	}
	defer response.Body.Close()
	//////////////////////////////////////// Proxy ////////////////////////////////////////////////////////////////

	// Read response content
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(errors.New("Fail to get the response"), "Response fail")
	}

	// Parsing JSON data
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, errors.Wrap(errors.New("Fail to parsing JSON data"), "Parse fail")
	}

	// 1. Check if the ABI exists in the database for the given chainID, contract address, and function signature.
	//    If found, return the ABI from the database

	// 2. If the ABi is not found in the database, query Etherscan to retrieve the ABI
	//    If Etherscan returns the ABI, store it in the database and return it

	// 3. If Etherscan does not have the ABI, return an appropriate error

	return []byte(apiResponse.Result), nil
}

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
