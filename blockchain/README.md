# fundslock
Smart contract and app for decentralized escrow system on blockchain

# Building

## Requirements

* Foundry - for basic SOlidity development and buliding
* Kurtosis - for integation tests running on local machine
* Python - for running static analysis using Slither
* make - for executing targets in Makefile

## Build project

Thr first time you check out the project, run the following command to initialize project 
dependencies and install required modules:

```
make init-project
```

To compile, test, format and lint project:

```
make build
```

To deploy contract on local Kurtosis node:

```
make deploy-contract-kurtosis
```

To run integration tests:

```
make intgration-tests
```

To perform static analysis:

```
make static-analysis
```

Test coverage report:

```
make coverage
```

View documentation:

```
make docs
```
## Deploying the contract

Details for deploying contract are specified in [deployment.md](./doc/deployment.md).

## Interacting with the contract

This section shows how to interact with the contract usinf `cast` tool. Requirements:

* cast (Foundry) installed
* Buyer address
* Seller address
* private key for each address
* some ETH to cover transaction fees and fund the agreement

All addresses can be created in MetaMask or similar wallet for testing purposes.

export the following environment variables:

* `CONTRACT_ADDRESS` - address of deployed contract
* `BUYER_ADDRESS` - buyer wallet address
* `SELLER_ADDRESS` - seller wallet address
* `BUYER_PRIVATE_KEY` - buyer wallet private key
* `SELLER_PRIVATE_KEY` - seller wallet private key
* `RPC_URL` - network RPC URL


### Create agreement

```
cast send $CONTRACT_ADDRESS "createAgreement(address,address payable,uint256)(uint256)" $SELLER_ADDRESS $BUYER_ADDRESS  10000000000000000 --rpc-url $RPC_URL --private-key $BUYER_PRIVATE_KEY
```

### Seler accept agreement

```
cast send $CONTRACT_ADDRESS "sellerAcceptAgreement(uint256)"  100 --rpc-url $RPC_URL --private-key $SELLER_PRIVATE_KEY
```
### Buyer fund agreement

```
cast send $CONTRACT_ADDRESS "fundAgreement(uint256)"  100 --value 0.01ether --rpc-url $RPC_URL --private-key $BUYER_PRIVATE_KEY
```

### Seller request release

```
cast send $CONTRACT_ADDRESS "requestRelease(uint256)"  100 --rpc-url $RPC_URL --private-key $SELLER_PRIVATE_KEY
```

### Release funds

```
cast send $CONTRACT_ADDRESS "releaseFunds(uint256)"  100  --rpc-url $RPC_URL --private-key $BUYER_PRIVATE_KEY
```