package main

import (
	"code/src/db"
	"code/src/fetch"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func main() {
	db := db.InitDatabase()

	fetcher := fetch.FetcherCli{"xxxx"}
	fetcher.GetABIAtStartOfBlock(db, 1, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), big.NewInt(12345))
}
