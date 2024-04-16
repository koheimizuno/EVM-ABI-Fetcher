package cache

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	addr := common.HexToAddress("0x123")
	abi := []byte("Hello, D23E:)")
	chainID := 1

	cache.Set(chainID, addr, abi)

	retrievedABI, found := cache.Get(chainID, addr)
	assert.True(t, found)
	assert.Equal(t, abi, retrievedABI)
}

// LRU logic: When the cache exceeds capacity, the least commonly used will be removed
func TestCacheEvictionPolicy(t *testing.T) {
	cache := NewABICache()
	// Fill cache beyond initial capacity
	for i := 0; i < 1001; i++ {
		addr := common.HexToAddress(fmt.Sprintf("0x%x", i))
		cache.Set(i, addr, []byte("TestCacheEvictionPolicy"))
	}

	// Originally supposed to be 0, the least commonly used,
	// After the query is completed, 0 is accessed, and the least commonly used becomes 1
	_, _ = cache.Get(0, common.HexToAddress("0x0"))

	// When a new data is inserted, 1 will be deleted
	cache.Set(1234, common.HexToAddress(fmt.Sprintf("0x%x", 1234)), []byte("TestCacheEvictionPolicy"))
	_, found := cache.Get(1, common.HexToAddress("0x1"))
	assert.False(t, found)
}

// Test the independent operation of the Evict method
func TestEvictFunction(t *testing.T) {
	cache := NewABICache()
	addr := common.HexToAddress("0x123")
	chainID := 1
	cache.Set(chainID, addr, []byte("ABI1"))
	cache.evict() // Force eviction regardless of cache size

	// Check if the element has been deleted
	_, found := cache.Get(chainID, addr)
	assert.False(t, found)
}
