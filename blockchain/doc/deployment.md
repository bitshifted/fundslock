
# Deploying project

Deploying the contract is done manually from the command line. This is done for security reasons to protect deployment keys and wallets.

Required software for deployment:
* UNIX-like system (Linux, Mac OS, Window WLS)
* Bash shell
* Foundry

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

### Create RPC node

Before trying to deploy the contract, you need an RPC node to communncate with. The easiest way is to signup with [Alchemy](https://www.alchemy.com/) (free tier is available), and create a project. This will give you RPC endpoints for chains where you want to deploy the contract.

You will need to specify endpoint for each chain you deploy to (ie. Ethereum, Arbitrum, Base etc.).

### Deploy contract

In addition to AWS environment variables defined above, define some additional variables:

* `AWS_KMS_KEY_ID` - ID of KMS key created in the first step
* `NETWORK_RPC_URL` - RPC URL of the network where you are deploying the contract
* `NETWORK_NAME` -  name of the network to deploy to. , Valid names are netowrk names from [Graph protocol](https://thegraph.com/docs/en/supported-networks/)

Then, run `make deploy-contract-public`. It will deploy the contract to requested network and print out contract address.

**Note:** You will need to repeat this process for each networn and RPC URL you want to support

