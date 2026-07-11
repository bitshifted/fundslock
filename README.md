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
