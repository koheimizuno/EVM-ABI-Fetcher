package cache

import (
	"container/list"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// CacheItem
// @dev Used to store a single cache entry in a linked list
type CacheItem struct {
	ChainID         int
	ContractAddress common.Address
	ABI             []byte
}

// ABICache
// @dev contain the cache strategy
type ABICache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

// NewABICache
// @dev Create a new ABICache with the given capacity
func NewABICache() *ABICache {
	cache := &ABICache{
		capacity: 1000, // Fixed capacity of 1000 entries
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
	return cache
}

// Get
// @dev Retrieve an item from the cache
func (c *ABICache) Get(chainID int, contractAddress common.Address) ([]byte, bool) {

	key := cacheKey(chainID, contractAddress)
	if element, found := c.cache[key]; found {
		c.list.MoveToFront(element)
		return element.Value.(*CacheItem).ABI, true
	}
	return nil, false
}

// Set
// @dev Add an item to the cache
func (c *ABICache) Set(chainID int, contractAddress common.Address, abi []byte) {
	key := cacheKey(chainID, contractAddress)

	newItem := &CacheItem{
		ChainID:         chainID,
		ContractAddress: contractAddress,
		ABI:             abi,
	}
	element := c.list.PushFront(newItem)
	c.cache[key] = element

	// Check capacity and remove the oldest items as necessary: LRU
	if c.list.Len() > c.capacity {
		c.evict()
	}

}

// evict LRU
func (c *ABICache) evict() {
	if element := c.list.Back(); element != nil {
		c.list.Remove(element)
		item := element.Value.(*CacheItem)
		delete(c.cache, cacheKey(item.ChainID, item.ContractAddress))
	}
}

// cacheKey
// @dev Generate key for cache mapping
func cacheKey(chainID int, address common.Address) string {
	return fmt.Sprintf("%d-%s", chainID, address.Hex())
}
