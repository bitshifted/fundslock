# FundsLock

FundsLock is  decentralized escrow system running on blockchain. It consists of smart contract, Graph protocol subgraph for events ang queries,  backend server and user-facing UI (web application). It runs on EVM compatible blockchains.

Projects comes with scripts and infrastructure provisioning modules to automate the complete deployment process. All modules are complete and ready for automated production deployment.

# Project structure

* [blockchain](./blockchain/) - Solidity/Foundry project implementing smart contract for escrow
* [backend](./backend/) - Go-based REST API server
* [graph](./graph/) - Graph protocol subgraph for events/log querying
* [frontend](./frontend/) - Frontend for the application

# Deploy on your own infrastructure

This sction describes how you can deploy the project on your own infrastructure . You will need access to the following services:

* Github - clone the repo and run automated deployment actions
* [Alchemy](https://alchemy.com) or [Infura](https://infura.io) - to run your own blockchain node
* AWS - for running frontend and backend service
* [Subgraph Studio](https://thegraph.com/studio/) - Graph Protocol UI for creating subgraphs

Required software and tools:

* Foundry
* Bash-compatible shell
* `make`
* `jq`

## Fork the repository

First step is to fork the repository to your own Github account. This allows you to run Github Actions with your own secrets and environment variables

## Deploy smart contract

Smart contract is deployesd to public blockchain manually rather than from Github Action. This is due to security reason and avoiding exposing wallet private keys. Follow the guide [here](./blockchain/doc/deployment.md).

## Deploy subgraphs

Smart contract deployment populates network information. Graphs are deployed automatically via Github Action. Complete process is described [here](./graph/README.md).
