package main

import (
	"code/src/db"
	"code/src/fetch"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load() // Load the `.env` file in the current directory by default
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apiKey := os.Getenv("API_KEY")
	rpcURL := os.Getenv("RPC_URL")

	db := db.InitDatabase()

	fetcher := fetch.FetcherCli{ApiKey: apiKey, RpcUrl: rpcURL}
	fmt.Println("First search")
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"))
	fmt.Println("Second search")
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"))
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"))
	_, err = fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("Successfully!")
	}
}
