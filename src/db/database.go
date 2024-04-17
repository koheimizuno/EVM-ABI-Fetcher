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
	SourceCode        string    `gorm:"type:text"`             // contract solidity source code TODO[修改了：增加]
	CompileTimeParams string    `gorm:"type:text"`             // constructor parameters TODO[修改了：增加]
	ContractABI       string    `gorm:"type:text"`             // The whole ABI of the contract  TODO[修改了：增加]
}

// FunctionSignature
// @dev Table 2
type FunctionSignature struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key"` // function signature unique identifier(uuid or int)
	ContractBytecodeID uuid.UUID `gorm:"type:uuid;index"`       // contract bytecode unique identifier (uuid or int) [foreign key] TODO[修改了：增加]
	Signature          []byte    `gorm:"type:blob;size:4"`      // function signature(hex or bytea, 4bytes)
	FunctionABI        string    `gorm:"type:text"`             // function ABI(json string)
}

// ContractDeployment
// @dev Table 3
type ContractDeployment struct {
	ChainID            int       `gorm:"type:int"`  // chainID(int)
	ContractAddress    []byte    `gorm:"type:blob"` // contract address(bytea or hex)
	ContractBytecodeID uuid.UUID `gorm:"type:uuid"` // contract bytecode unique identifier(uuid or int)
	//FunctionSignatureID uuid.UUID `gorm:"type:uuid"` // function signature unique identifier(uuid o int) TODO[修改了: 删除]
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
		!db.Migrator().HasTable(&ContractDeployment{}) {
		db.AutoMigrate(&ContractBytecode{}, &FunctionSignature{}, &ContractDeployment{})
		fmt.Println("Init the data successfully!")
	}

	return db
}
