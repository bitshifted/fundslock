
# Graph protocol setup

This component creates and deploys Graph protocol subgraph for FundsLock smart contract. It enables querying events and 
changes emitted from the contract.

Subgraphs are supported on all networks where contract is deployed.


## Setup 

Subgraph must be created in [Subgraph Studio](https://thegraph.com/studio/), with the following slug format per network: `fundsloack-<network name>`. For example, following are valid fomats:

* Ethereum mainnet: `fundslock-ethereum`
* Sepolia testnet: `fundslock-sepolia`
* Arbitrum mainnet: `fundslock-arbitrum`
* Arbitrum Sepolia tesnet: `fundslock-arbitrum-sepolia`

Once the subgraph is created, note the API key associated with it (for each graph).

## Deployment

Before attempting deployment, create Github secret for each of the graphs you deploy. Secret name should follow format `GRAPH_<network name>_KEY`, so it can be automatically used in Github Action deployment.

For example, for graphs created in previous section, secret names should be:

* `GRAPH_ETHEREUM_KEY`
* `GRAPH_SEPOLIA_KEY`
* `GRAPH_ARBITRUM_KEY`
* `GRAPH_ARBITRUM_SEPOLIA_KEY`



Deployment is done autmomatically via Github Action. Process assumes that smart contract is alrteady deployed to each network following the instructions provided in [blockchain documentation](../blockchain/doc/deployment.md):

1. Smart contract deployment will generate subgraph deployment information in this directory, specifically `networks.json` file. This is used to configure deployment information.
2. create a branch and push changes to Github. This will force buinding and testing your subgraph.
3. Once all tests pass, merge the pull request. This will trigger Github workflow to deploy the subgraph

## Publishing the subgraph

Deployment process only deploys subgraphs to target network, and they can be accessed via rate limited endpoints. For production deployments, subgraphs need to be published.

Publishing requires having crypto wallet with GRT token with which you pay the publishing. For security reasons and fact that this involves real-world money, publishing is not automated, but rather should be done via Subgraph studion UI.


# Troubleshooting

In test/dev environemnt, if graph is not queried for a prolonged time, it may be purged. The solution is to force a dummy commit into this directory, triggering a fresh release to dev/test environment.

One symptom of this happening is if you get this kind of response when querying graph:

```
{"errors":[{"message":"deployment `u1756772/s120401/latest` does not exist"}]}
```

