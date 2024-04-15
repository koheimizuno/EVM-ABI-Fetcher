package main

import (
	"code/src/db"
	"code/src/fetch"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"math/big"
	"os"
)

func main() {
	err := godotenv.Load() // Load the `.env` file in the current directory by default
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apiKey := os.Getenv("API_KEY")

	db := db.InitDatabase()

	fetcher := fetch.FetcherCli{apiKey}
	data, err := fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), big.NewInt(12345))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("data", string(data))
	}
}
