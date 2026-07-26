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

# Deploying project

## Deploying to local Kurtosis node

Project can spin up local Kurtosis node for testing. This is done automatically by integration tests, but can alsso be used for local contract deployment, for testing purposes.

```
make deploy-contract-kurtosis
```

This target will spin up local node and deploy contract to it.

## Deploying to test/mainnet using AWS KMS

Deployment to testnet or mainnet is supported by using AWS KMS for private key storage.

### Import private key to AWS KMS

First step is to create KMS key with external key material. Wallet private key for deployment will be uploaded in it.

The following environment variables need to be set:
* `AWS_REGION`
* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`

Run the script `infra/kms-key-create.sh --policy-path my-kms-key-policy.json`. You need to specify path to KMS key policy path, so that it can be applied to the key. Example policy can be found in `infra/kms-key-policy.json`.

This will create KMS key in AWS. The script returns ID of the key thet was created. Check your AWS console for the newly created key.

### Import key material from wallet private key

Export private key that will be used for deployment from your wallet. Make sure that no leading `0x` is present.

Run the script `infra/secure-kms-import.sh <kms-key-id>`,  where `kms-key-id` is ID of the key created in previous step. The script will prompt you to enter wallet private key. Simply paste it in terminal and press `Enter`. Script will perform key encryption and upload key material to KMS. 

The script creates secure RAM disk for the keys and performs cleanup after that, so that no trace of the keys is left on the computer.

### Deploy contract

In addition to AWS environment variables defined above, define some additional variables:

* `AWS_KMS_KEY_ID` - ID of KMS key created in the first step
* `NETWORK_RPC_URL` - RPC URL of the network where you are deploying the contract
* `NETWORK_NAME` - (optional) name of the network to deploy to. Default is `sepolia`, whihc deploys to sepolia testnet, Valid names are netowrk names from [Graph protocol](https://thegraph.com/docs/en/supported-networks/)

Then, run `make deploy-contract-public`. It will deploy the contract to requested network and print out contract address.

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