package cache

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

var signature = "0xa9059cbb"
var functionABI = "{ \"constant\": false, \"inputs\": [], \"name\": \"world\", \"outputs\": [ { \"name\": \"\", \"type\": \"uint256\" } ], \"payable\": true, \"stateMutability\": \"payable\", \"type\": \"function\" }"
var contractABI = "[\n\t{\n\t\t\"constant\": true,\n\t\t\"inputs\": [],\n\t\t\"name\": \"hello\",\n\t\t\"outputs\": [\n\t\t\t{\n\t\t\t\t\"name\": \"\",\n\t\t\t\t\"type\": \"string\"\n\t\t\t}\n\t\t],\n\t\t\"payable\": false,\n\t\t\"stateMutability\": \"view\",\n\t\t\"type\": \"function\"\n\t},\n\t{\n\t\t\"constant\": false,\n\t\t\"inputs\": [],\n\t\t\"name\": \"world\",\n\t\t\"outputs\": [\n\t\t\t{\n\t\t\t\t\"name\": \"\",\n\t\t\t\t\"type\": \"uint256\"\n\t\t\t}\n\t\t],\n\t\t\"payable\": true,\n\t\t\"stateMutability\": \"payable\",\n\t\t\"type\": \"function\"\n\t}\n]"
var chainID = 1
var contractAddress = common.HexToAddress("0x67804f76158410B1C5aE84fED4220E7bf5c1F9dE")

// To ensure that the cache is properly initialized
func TestNewABICache(t *testing.T) {
	cache := NewABICache()
	assert.Equal(t, 1000, cache.capacity)
	assert.NotNil(t, cache.cache)
	assert.NotNil(t, cache.list)
}

// The normal operation of the Set and Get methods
func TestCacheSetAndGet(t *testing.T) {
	cache := NewABICache()
	cache.Set(chainID, contractAddress, []byte(functionABI), []byte(contractABI), signature)

	FunctionABIGet, ContractABIGet, isFound := cache.Get(chainID, contractAddress, signature)
	assert.True(t, isFound)
	assert.Equal(t, FunctionABIGet, []byte(functionABI))
	assert.Equal(t, ContractABIGet, []byte(contractABI))
}

// LRU logic: When the cache exceeds capacity, the least commonly used will be removed
func TestCacheEvictionPolicy(t *testing.T) {
	cache := NewABICache()
	// Fill cache beyond initial capacity
	for i := 0; i < 1001; i++ {
		addr := common.HexToAddress(fmt.Sprintf("0x%x", i))
		cache.Set(i, addr, []byte(functionABI), []byte(contractABI), signature)
	}

	// Originally supposed to be 0, the least commonly used,
	// After the query is completed, 0 is accessed, and the least commonly used becomes 1
	_, _, _ = cache.Get(0, common.HexToAddress("0x0"), signature)

	// When a new data is inserted, 1 will be deleted
	cache.Set(1234, contractAddress, []byte("123"), []byte("456"), signature)
	_, _, found := cache.Get(1, common.HexToAddress("0x1"), signature)
	assert.False(t, found)
}

// Test the independent operation of the Evict method
func TestEvictFunction(t *testing.T) {
	cache := NewABICache()

	cache.Set(chainID, contractAddress, []byte(functionABI), []byte(contractABI), signature)
	cache.evict() // Force eviction regardless of cache size

	// Check if the element has been deleted
	_, _, found := cache.Get(chainID, contractAddress, signature)
	assert.False(t, found)
}
