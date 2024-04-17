package cache

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Prepare some data that may be used
var (
	contractAddress = common.HexToAddress("0xc0ffee254729296a45a3885639AC7E10F9d54979")
	functionABI     = &abi.Method{Name: "transfer", RawName: "transfer(address,uint256)"}
	contractABI     = &abi.ABI{Methods: map[string]abi.Method{}}
	signature       = "transfer(address,uint256)"
)

func TestNewABICache(t *testing.T) {
	cache := NewABICache()
	assert.NotNil(t, cache)
	assert.Equal(t, 1000, cache.capacity)
	assert.NotNil(t, cache.list)
	assert.NotNil(t, cache.cache)
}

func TestSetAndGetCacheItem(t *testing.T) {
	cache := NewABICache()
	cache.Set(1, contractAddress, functionABI, contractABI, signature)

	// test Getter
	fetchedMethod, fetchedAbi, found := cache.Get(1, contractAddress, signature)
	assert.True(t, found)
	assert.Equal(t, functionABI, fetchedMethod)
	assert.Equal(t, contractABI, fetchedAbi)

	// Is the test item at the forefront of the cache
	assert.Equal(t, cache.list.Front().Value.(*CacheItem).Signature, signature)
}

func TestCacheEviction(t *testing.T) {
	cache := NewABICache()
	cache.capacity = 3 // Set small capacity for easy testing

	// Fill cache
	cache.Set(1, contractAddress, functionABI, contractABI, signature)
	cache.Set(2, contractAddress, functionABI, contractABI, "function2")
	cache.Set(3, contractAddress, functionABI, contractABI, "function3")

	_, _, _ = cache.Get(1, contractAddress, signature) // the second Item will be the RLU

	// Add the forth item to trigger evict()
	cache.Set(4, contractAddress, functionABI, contractABI, "function4")

	_, _, found := cache.Get(1, contractAddress, signature)
	assert.False(t, found) // the second Item will be the RLU

	_, _, found = cache.Get(2, contractAddress, "function2")
	assert.True(t, found) // the first item is still exist

	_, _, found = cache.Get(3, contractAddress, "function3")
	assert.True(t, found) // the third item is still exist
}
