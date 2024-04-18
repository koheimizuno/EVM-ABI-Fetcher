package cache

import (
	"container/list"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// CacheItem
// @dev Used to store a single cache entry in a linked list
type CacheItem struct {
	ChainID         int
	ContractAddress common.Address
	Signature       string      // E.g. transfer(address,uint256) => we store 0xa9059cbb
	FunctionABI     *abi.Method // the ABI of the Signature.
	ContractABI     *abi.ABI    // The whole ABI of the contract
}

// ABICache
// @dev contain the cache strategy
type ABICache struct {
	capacity int
	cache    map[int64]*list.Element
	list     *list.List
}

// NewABICache
// @dev Create a new ABICache with the given capacity
func NewABICache() *ABICache {
	cache := &ABICache{
		capacity: 1000, // Fixed capacity of 1000 entries
		cache:    make(map[int64]*list.Element),
		list:     list.New(),
	}
	return cache
}

// Get
// @dev Retrieve an item from the cache
// @return FunctionABI, ContractABI, isFound
// Notice: chainID+contractAddress+signature => return functionABI
// Notice: chainID+contractAddress+"Search for contractABI" => return contractABI
func (c *ABICache) Get(chainID int, contractAddress common.Address, signature string) (functionABI *abi.Method, contractABI *abi.ABI, isFound bool) {
	key := CacheKey(chainID, contractAddress, signature)
	if element, found := c.cache[key]; found {
		c.list.MoveToFront(element)
		return element.Value.(*CacheItem).FunctionABI, element.Value.(*CacheItem).ContractABI, true
	}
	return nil, nil, false
}

// Set
// @dev Add an item to the cache
func (c *ABICache) Set(chainID int, contractAddress common.Address, functionABI *abi.Method, contractABI *abi.ABI, signature string) {
	key := CacheKey(chainID, contractAddress, signature)

	newItem := &CacheItem{
		ChainID:         chainID,
		ContractAddress: contractAddress,
		Signature:       signature,
		FunctionABI:     functionABI,
		ContractABI:     contractABI,
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
		delete(c.cache, CacheKey(item.ChainID, item.ContractAddress, item.Signature))
	}
}

// cacheKey
// @dev Generate key for cache mapping
// To find functionABI: signature => signature
// To find contractABI: signature => ""
func CacheKey(chainID int, address common.Address, signature string) int64 {

	input := fmt.Sprintf("%d-%s-%s", chainID, address.Hex(), signature)

	hash := crypto.Keccak256([]byte(input))

	// high 8 bytes
	high8Bytes := hash[len(hash)-8:]

	// bytes => big.Int => int64
	high8BigInt := new(big.Int).SetBytes(high8Bytes)
	high8Int64 := high8BigInt.Int64()

	return high8Int64
}

// Get4bytesSig 使用以太坊的Keccak256哈希算法来生成输入字符串的哈希值，并返回最高4字节
func Get4bytesSig(signature string) [4]byte {
	hash := crypto.Keccak256([]byte(signature)) // 使用Keccak256哈希算法

	// 取最高的4字节
	var result [4]byte
	copy(result[:], hash[len(hash)-4:])

	return result
}
