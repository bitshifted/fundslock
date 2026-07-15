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

# Deploying project

## Deploying to local Kurtosis node

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

Run the script `infra/kms-key-import.sh <kms-key-id>`,  where `kms-key-id` is ID of the key created in previous step. The script will prompt you to enter wallet private key. Simply paste it in terminal and press `Enter`. Script will perform key encryption and upload key material to KMS. 

The script creates secure RAM disk for the keys and performs cleanup after that, so that no trace of the keys is left on the computer.

### Deploy contract

In addition to AWS environment variables defined above, define some additional variables:

* `AWS_KMS_KEY_ID` - ID of KMS key created in the first step
* `NETWORK_RPC_URL` - RPC URL of the network where you are deploying the contract

Then, run `make deploy-contract-public`. It will deploy the contract to requested network and print out contract address.
