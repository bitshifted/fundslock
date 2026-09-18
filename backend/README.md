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

### Configuration variables

Copy file `.env.example` as `.env`. The following environment variables can be set:

* `JWT_SECRET_KEY` (required) - secret key used for JWT token signing
* `COOKIE_DOMAIN` (required) - sets domain for cookies
* `GRAPH_CONFIG_BASE64` (required) - access data for Graph protocol subgraphs. This is Base64 encoded JSON string. Example:

```shell
echo -n '[{"chain":"sepolia","graphUrl":"https://somehost.com/graph","apiKey":"abcdef"}]' | base64 -w 0
# W3siY2hhaW4iOiJzZXBvbGlhIiwiZ3JhcGhVcmwiOiJodHRwczovL3NvbWVob3N0LmNvbS9ncmFwaCIsImFwaUtleSI6ImFiY2RlZiJ9XQ==
GRAPH_CONFIG_BASE64=W3siY2hhaW4iOiJzZXBvbGlhIiwiZ3JhcGhVcmwiOiJodHRwczovL3NvbWVob3N0LmNvbS9ncmFwaCIsImFwaUtleSI6ImFiY2RlZiJ9XQ==
```
* `ACCESS_TOKEN_DURATION` (optional) - how long is acces token valid, in seconds. Default value is 300 (5 minutes)
* `REFRESH_TOKEN_DURATION` (otpional) - how long refresh token is valid, in seconds. Defaukt value is 3600000 (7 days)
* `SECCURE_COOKIE` (optional) - whether to use secure cookies. Defaults to `true`

**Note:** Using `.env` files in production is discouraged for security reasons. Only use it for development purposes.


# Authentication

Backend uses JWT tokens and wallet signing to authenticate users. Once user connetcs the wallet in frontend, they are prompted to authenticate with SIWE (Sign In With Ethereum) message. Based on the message, backend issues access token and refresh token.


