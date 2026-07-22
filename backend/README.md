# Backend

REST API backend for FundsLock smart contract. This is Go project that runs a HTTP server and has the following functions:

* authentication
* communication with Graph protocol

# Building and running

Requirements:
* Go version 1.25.x
* make
* Docker (optional)

To build the project, run

```
make build
```

This performs linting, building and testing, including coverage report. Build artifacts can be found under `target/<platform>`, where each `<platform>` directory contains binary for the given platform. 

To run the project locally, use `go run main.go start`. This start the server on port 3000.

