package db

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @dev Table 1
type ContractBytecode struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key"` // contract bytecode unique identifier(uuid or int)
	Bytecode []byte    `gorm:"type:blob"`             // contract bytecode(hex or bytea)
}

// @dev Table 2
type FunctionSignature struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"` // function signature unique identifier(uuid or int)
	Signature   []byte    `gorm:"type:blob;size:4"`      // function signature(hex or bytea, 4bytes)
	FunctionABI string    `gorm:"type:text"`             // function ABI(json string)
}

// @dev Table 3
type ContractDeployment struct {
	ChainID               int       `gorm:"type:int"`  // chainID(int)
	ContractAddress       []byte    `gorm:"type:blob"` // contract address(bytea or hex)
	ContractBytecodeID    uuid.UUID `gorm:"type:uuid"` // contract bytecode unique identifier(uuid or int)
	FunctionSignatureID   uuid.UUID `gorm:"type:uuid"` // function signature unique identifier(uuid o int)
	DeployedBlockNumber   int       `gorm:"type:int"`  // deployedAt - block number
	DeployedTxIndex       int       `gorm:"type:int"`  // deployedAt - txIndex
	DestructedBlockNumber int       `gorm:"type:int"`  // destructedAt - block number
	DestructedTxIndex     int       `gorm:"type:int"`  // destructedAt - txIndex
}

// @dev Init the database, create 3 tables
func InitDatabase() (db *gorm.DB) {
	// Open the database in the current directory, create a new one if it is not exist
	db, err := gorm.Open(sqlite.Open("ABIs.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, please check the SQLite3 environment")
	}

	// Create 3 tables
	db.AutoMigrate(&ContractBytecode{}, &FunctionSignature{}, &ContractDeployment{})

	return db
}
