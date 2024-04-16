package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test database initialization
func TestInitDatabase(t *testing.T) {
	db := InitDatabase()
	if db == nil {
		t.Errorf("TestInitDatabase fail")
	}

	// Check if a table has been created
	assert.True(t, db.Migrator().HasTable(&ContractBytecode{}))
	assert.True(t, db.Migrator().HasTable(&FunctionSignature{}))
	assert.True(t, db.Migrator().HasTable(&ContractDeployment{}))
}
