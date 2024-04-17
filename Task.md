# ABI Caching Project

> The `README.md` will be updated when the project is completed

![image-20240415085028979](Task/image-20240415085028979.png)

## Objective

The goal of this project is to cache smart contract ABIs in a database and provide a function to retrieve the ABI given a chain ID, contract address, and block number. The function should perform lazy fetching and caching to optimize performance.

## Database Schema

Create a database table using GORM and SQLite3 with the following schema:

- [ ] Table 1:
  - [x] contract bytecode unique identifier(uuid or int) [primary key]
  - [x] contract bytecode(hex or bytea)
  - [ ] contract source code ← (work on this only if you finished all other tasks)
  - [ ] parameters at compile time (e.g., how many rounds of optimisation) ← (work on this only if you finished all other tasks)
- [ ] Table 2:
  - [x] function signature unique identifier(uuid or int)
  - [ ] contract bytecode unique identifier (uuid or int) [foreign key]
  - [ ] function signature(hex or bytea, 4bytes)
  - [ ] function ABI(json string)
- [x] Table 3:
  - [x] chainID(int)
  - [x] contract address(bytea or hex)
  - [x] contract bytecode unique identifier(uuid or int)
  - [ ] ~~deployedAt <= (work on this only if you finished all other tasks)~~
    - [ ] ~~block number~~
    - [ ] ~~txIndex~~
  - [ ] ~~destructedAt <= (work on this only if you finished all other tasks)~~
    - [ ] ~~block number~~
    - [ ] ~~txIndex~~

## Get ABI Function (In GoLand)

- [x] Implement a function with the following signature:

```go
func GetABIAtStartOfBlock(chainID int, contractAddress common.Address, block *big.Int) ([]byte, error)
=========>
func GetABIAtStartOfBlock(chainID int, contractAddress common.Address) ([]byte, error)
```

The function should perform the following steps:

- [x] Check if the ABI exists in the database for the given chainID, contract address, and function signature.

  > If found, return the ABI from the database

- [x] If the ABi is not found in the database, query Etherscan to retrieve the ABI

  > If Etherscan returns the ABI, store it in the database and return it

- [x] If Etherscan does not have the ABI, return an appropriate error

## Caching

- [x] Implement an in-memory cache using a map to store the most recently queried ABIs.
- [x] Use a cache size of 1000 entries
- [x] Implement a least recently used(LRU) eviction policy to remove the least recently accessed entries when the cache reaches its maximum size.
- [x] Ensure thread-safety for concurrent access to the cache using a sync.RWMutex

## Error Handing and Logging

- [x] If there is a timeout on Etherscan, wait and retry
- [x] Handle errors gracefully and return appropriate error messages from the GetABI function
- [x] Log errors and key events using the logrus logging package with the following log levels:
  - [x] Error: For critical errors that prevent the function from executing properly
  - [x] Warning: For non-critical issues or unexpected behavior
  - [x] Info: For important events or milestones during the execution
  - [x] Log the input parameters, retrieved ABI, and any error messages for debugging purposes

## Performance

- [ ] Optimize database queries by creating appropriate indexes on the ChainID, ContractAddress, and FuncSignature columns using GORM: Re indexes, we are okay with slow inserts, but we want very fast query speed. Do you create indexes for your tables?
- [x] Implement caching to minimize the number of database queries and external API calls
- [x] Aim for a maximum response time of 100ms for the GetABI function

## Testing and Validation

- [ ] Unit Tests:
  - [x] Write unit test using the `testing` package in Go to cover the core functionality of the GetABI function
  - [ ] Test scenarios where the ABI is found in the database and retrieved from Etherscan. (Lets imagine later, we encountered a transaction, which calls into a contract without ABI. We know the data is “0xAAAA…..“. Now what we can do is, I would look up in table 2, and find all abis matches to 0xAAAA, and try each to see if it can decode the data.)
  - [x] Test error handling for scenarios where the ABI is not found or Etherscan returns an error
- [ ] Integration Tests:
  - [x] Write integration test to verify the interaction between the GetABI function, the database, and Etherscan
  - [x] Test the caching mechanism to ensure that frequently accessed ABIs are retrieved from the cache
  - [ ] ~~Ensure that the proxy contract table is properly populated during the tests~~
- [x] Validation:
  - [x] Compare the retrieved ABIs with the expected ABIs from reliable sources to ensure correctness
  - [x] Handle any discrepancies or inconsistencies between the retrieved ABIs and the expected ABIs

## Deliverables

- [ ] Pull Request on Github
- [ ] Go source code for the GetABI function and associated helper functions
- [x] Unit and interation test suite using the `testing` package
- [x] Documentation(README.md) explaining the design choices, assumptions, and any dependencies

## Timeline

- [x] Day 1: Set up the development environment, create the database schema using GORM, and implement the basic GetABI function
- [x] Day 2: integrate Etherscan for ABI retrieval and implement caching using a map and sync.RWMutex
- [ ] Day 3-4: Conduct through testing validation, and performance optimization
- [ ] Day 5: Prepare documentation, conduct final code review, and address any remaining issues

