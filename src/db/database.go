package db

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ContractBytecode
// @dev Table 1
type ContractBytecode struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key"` // contract bytecode unique identifier(uuid or int)
	Bytecode []byte    `gorm:"type:blob"`             // contract bytecode(hex or bytea)
}

// FunctionSignature
// @dev Table 2
type FunctionSignature struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"` // function signature unique identifier(uuid or int)
	Signature   []byte    `gorm:"type:blob;size:4"`      // function signature(hex or bytea, 4bytes)
	FunctionABI string    `gorm:"type:text"`             // function ABI(json string)
}

// ContractDeployment
// @dev Table 3
type ContractDeployment struct {
	ChainID             int       `gorm:"type:int"`  // chainID(int)
	ContractAddress     []byte    `gorm:"type:blob"` // contract address(bytea or hex)
	ContractBytecodeID  uuid.UUID `gorm:"type:uuid"` // contract bytecode unique identifier(uuid or int)
	FunctionSignatureID uuid.UUID `gorm:"type:uuid"` // function signature unique identifier(uuid o int)
	// The field below will not be used temporarily
	DeployedBlockNumber   int `gorm:"type:int"` // deployedAt - block number
	DeployedTxIndex       int `gorm:"type:int"` // deployedAt - txIndex
	DestructedBlockNumber int `gorm:"type:int"` // destructedAt - block number
	DestructedTxIndex     int `gorm:"type:int"` // destructedAt - txIndex
}

// InitDatabase
// @dev Init the database, get the database's handle
// @return SQLite3's handle
func InitDatabase() (db *gorm.DB) {
	// Get the database's handle
	db, err := gorm.Open(sqlite.Open("ABIs.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	// Check if tables exist and migrate if they do not
	if !db.Migrator().HasTable(&ContractBytecode{}) ||
		!db.Migrator().HasTable(&FunctionSignature{}) ||
		!db.Migrator().HasTable(&ContractDeployment{}) {
		db.AutoMigrate(&ContractBytecode{}, &FunctionSignature{}, &ContractDeployment{})
		fmt.Println("Init the data successfully!")
	}

	return db
}
