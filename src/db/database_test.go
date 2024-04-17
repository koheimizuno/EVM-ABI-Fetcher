package db

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var db = InitDatabase()

// Test database initialization
func setup(t *testing.T) {
	if db == nil {
		t.Errorf("TestInitDatabase fail")
	}

	// Check if a table has been created
	assert.True(t, db.Migrator().HasTable(&ContractBytecode{}))
	assert.True(t, db.Migrator().HasTable(&FunctionSignature{}))
	assert.True(t, db.Migrator().HasTable(&SearchEtherscan{}))
	assert.True(t, db.Migrator().HasTable(&ContractDeployment{}))
}

func tearDown() {
	db.Exec("DELETE FROM contract_bytecodes")
	db.Exec("DELETE FROM function_signatures")
	db.Exec("DELETE FROM contract_deployments")
	db.Exec("DELETE FROM search_etherscans")
}

func TestContractBytecode(t *testing.T) {
	setup(t)
	defer tearDown()

	cb := ContractBytecode{
		ID:                uuid.New(),
		Bytecode:          []byte{0xab, 0xcd, 0xef},
		SourceCode:        "pragma solidity ^0.8.0;",
		CompileTimeParams: "param1,param2",
		ContractABI:       "ABI details here",
	}
	result := db.Create(&cb)
	assert.Nil(t, result.Error)
}

func TestSearchEtherscan(t *testing.T) {
	setup(t)
	defer tearDown()

	cb := SearchEtherscan{
		ChainID:         123,
		ContractAddress: []byte("123"),
		Time:            123123,
		ShouldSearch:    true,
	}
	result := db.Create(&cb)
	assert.Nil(t, result.Error)
}

func TestFunctionSignature(t *testing.T) {
	setup(t)
	defer tearDown()

	cb := ContractBytecode{
		ID:                uuid.New(),
		Bytecode:          []byte{0xab, 0xcd, 0xef},
		SourceCode:        "pragma solidity ^0.8.0;",
		CompileTimeParams: "param1,param2",
		ContractABI:       "ABI details here",
	}
	db.Create(&cb)

	fs := FunctionSignature{
		ID:                 123456789,
		ContractBytecodeID: cb.ID,
		Signature:          []byte{0x1a, 0x2b, 0x3c, 0x4d},
		FunctionABI:        "{\"inputs\": [], \"name\": \"testFunc\", \"type\": \"function\"}",
	}
	result := db.Create(&fs)
	assert.Nil(t, result.Error)
}

func TestContractDeployment(t *testing.T) {
	setup(t)
	defer tearDown()

	cb := ContractBytecode{
		ID:                uuid.New(),
		Bytecode:          []byte{0xab, 0xcd, 0xef},
		SourceCode:        "pragma solidity ^0.8.0;",
		CompileTimeParams: "param1,param2",
		ContractABI:       "ABI details here",
	}
	db.Create(&cb)

	cd := ContractDeployment{
		ChainID:            1,
		ContractAddress:    []byte{0x00, 0x1a, 0x2b, 0x3c, 0x4d},
		ContractBytecodeID: cb.ID,
	}
	result := db.Create(&cd)
	assert.Nil(t, result.Error)
}
