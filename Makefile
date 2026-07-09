
SHELL := /bin/bash

.PHONY: check-python

check-python:
	@command -v python3 >/dev/null 2>&1 || { echo "Error: Python3 is not installed."; exit 1; }
	@echo "Python3 is installed at: $$(which python3)"


# Initialize project (create virtual environment and install dependencies)
init-project: check-python
	@echo "Initializing project..."
	@python3 -m venv .venv
	@echo "Installing dependencies..."
	@.venv/bin/pip3 install -r requirements.txt
	@forge install

copyright-check:
	@echo "Checking copyright headers..."
	@grep -Lr "^// Copyright (c) 2026 Bitshift ED$$" src/*.sol test/*.sol script/*.sol


format-check:
	@echo "Checking code formatting..."
	@forge fmt --check

lint:
	@echo "Running linter..."
	@forge lint

test:
	@echo "Running tests..."
	@forge test

build: format-check copyright-check lint test
	@echo "Building project..."
	@forge build

start-testnet:
	@running=$$(kurtosis enclave inspect local-eth-testnet | grep "Status:" | awk '{print $$NF}'); \
	if [ "$$running" = "RUNNING" ]; then \
			echo "Kurtosis enclave 'local-eth-testnet' is already running."; \
	else \
			echo "Starting local Ethereum testnet using Kurtosis..."; \
			kurtosis --enclave local-eth-testnet run github.com/ethpandaops/ethereum-package; \
	fi

deploy-contract-kurtosis: start-testnet
	@echo "Deploying BasicEscrow contract..."; \
	KURTOSIS_RPC_URL=$$(kurtosis port print local-eth-testnet el-1-geth-lighthouse rpc); \
	forge script script/BasicEscrowDeploy.s.sol --rpc-url $$KURTOSIS_RPC_URL --broadcast

integration-test: deploy-contract-kurtosis
	@echo "Running integration tests..."; \
	KURTOSIS_RPC_URL=$$(kurtosis port print local-eth-testnet el-1-geth-lighthouse rpc); \
	forge script script/IntegrationTest.s.sol --rpc-url $$KURTOSIS_RPC_URL --broadcast

make static-analysis:
	@echo "Running static analysis..."; 
	@.venv/bin/slither . --solc-remaps "@openzeppelin/=lib/openzeppelin-contracts/"
