package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ContractBytecode
// @dev Table 1
type ContractBytecode struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key"` // contract bytecode unique identifier(uuid or int)
	Bytecode          []byte    `gorm:"type:blob"`             // contract bytecode(hex or bytea)
	SourceCode        string    `gorm:"type:text"`             // contract solidity source code
	CompileTimeParams string    `gorm:"type:text"`             // constructor parameters
	ContractABI       string    `gorm:"type:text"`             // The whole ABI of the contract
}

// FunctionSignature
// @dev Table 2
type FunctionSignature struct {
	ID                 int64     `gorm:"type:bigint;primary_key"` // function signature unique identifier(int, 8 bytes)
	ContractBytecodeID uuid.UUID `gorm:"type:uuid;index"`         // contract bytecode unique identifier (uuid or int) [foreign key]
	Signature          []byte    `gorm:"type:blob;size:4"`        // function signature(hex or bytea, 4bytes)
	FunctionABI        string    `gorm:"type:text"`               // function ABI(json string)
}

// ContractDeployment
// @dev Table 3
type ContractDeployment struct {
	ChainID            int       `gorm:"type:int"`  // chainID(int)
	ContractAddress    []byte    `gorm:"type:blob"` // contract address(bytea or hex)
	ContractBytecodeID uuid.UUID `gorm:"type:uuid"` // contract bytecode unique identifier(uuid or int)
}

// SearchEtherscan represents a table structure for blockchain scanning options
type SearchEtherscan struct {
	ChainID         int    `gorm:"type:int;index"` // Chain ID as integer
	ContractAddress []byte `gorm:"type:blob"`      // Contract address in byte array or hex
	Time            int    `gorm:"type:int"`       // Time as integer (e.g., UNIX timestamp)
	ShouldSearch    bool   `gorm:"type:boolean"`   // Flag to indicate if a search should be performed
}

var log = logrus.New()

// InitDatabase
// @dev Init the database, get the database's handle
// @return SQLite3's handle
func InitDatabase() (db *gorm.DB) {
	// Get the database's handle
	db, err := gorm.Open(sqlite.Open("ABIs.db"), &gorm.Config{})
	if err != nil {
		log.Error("Fail to connect to the database in current directory: ABIs.db")
		panic("Fail to connect to the database in current directory: ABIs.db")
	}

	// Check if tables exist and migrate if they do not
	if !db.Migrator().HasTable(&ContractBytecode{}) ||
		!db.Migrator().HasTable(&FunctionSignature{}) ||
		!db.Migrator().HasTable(&SearchEtherscan{}) ||
		!db.Migrator().HasTable(&ContractDeployment{}) {
		db.AutoMigrate(&ContractBytecode{}, &FunctionSignature{}, &ContractDeployment{}, &SearchEtherscan{})
		fmt.Println("Init the data successfully!")
	}

	return db
}
