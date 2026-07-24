
# Graph protocol setup

This component creates and deploys Graph protocol subgraph for FundsLock smart contract. It enables querying events and 
changes emitted from the contract.

Currently supports only subgraphs on Sepolia testnet.


## Setup 

Subgraph must be created in [Subgraph Studio](https://thegraph.com/studio/), with the following slugs per network:

* Sepoliaa testner: `fundslock-sepolia-eth`

Once the subgraph is created, note the API key associated with it.

## Deployment

Deployment is done autmomatically via Github Action. The process is the following:

1. dpeloy smart contract using instructions provided in [blockchain README](../blockchain/README.md#deploying-to-testmainnet-using-aws-kms)
2. Deployment will generate subgraph deployment information in this directory, specifically `networks.json` file. This is used to configure deployment information.
3. create a branch and push changes to Github. This will force buinding and testing your subgraph.
4. Once all tests pass, merge the pull request. This will trigger Github workflow to deploy the subgraph



