# Agent Guidelines & Runbook (AGENTS.md)

Welcome, AI Agent! This file serves as your comprehensive instruction manual and onboarding guide for working in the **blockchain** directory. It contains our architectural layout, domain model, operational commands, coding guidelines, and testing requirements to ensure your contributions align perfectly with our standards.

---

## 1. Repository Overview

**fundslock** is a decentralized escrow system implemented as an EVM smart contract, with off-chain backend and frontend component. It enables buyers and sellers to establish secure agreements, fund contracts, request/approve releases, and engage in mediator-driven dispute resolutions. This directory (`blockchain`) contains smart contract implementation.

### Key Technology Stack
- **Smart Contracts:** Solidity (`^0.8.35`)
- **Framework & Tooling:** Foundry (`forge` / `cast`)
- **Integration Testing:** Kurtosis (orchestrates ephemeral multi-container Ethereum testnets)
- **Security & Analysis:** Slither (static analysis tool via a local Python 3 environment)
- **Task Orchestration:** `Makefile`

---

## 2. Directory Architecture

```
/
├── config/                     # Configuration files for test accounts
├── infra/                      # Scripts and IAM policies for secure AWS KMS key imports
│   ├── kms-key-create.sh       # Sinks a KMS key configured with custom policies
│   ├── kms-key-policy.json     # Policy definition for the KMS deployment key
│   └── secure-kms-import.sh    # Securely mounts a RAM disk to upload keys to AWS KMS
├── lib/                        # Git submodules for external libraries (OpenZeppelin, forge-std)
├── script/                     # Deployment and integration scripts
│   ├── common/
│   │   └── Wallets.s.sol       # Helper script defining wallet structs/keys
│   ├── deploy/
│   │   ├── DeployLocal.s.sol   # Script for local Kurtosis deployments
│   │   └── DeployPublic.s.sol  # Script for KMS-backed public network deployments
│   └── IntegrationTest.s.sol   # Full integration test scenario suite
├── src/                        # Smart Contract Source Code
│   ├── Errors.sol              # Repository custom errors definitions
│   ├── FundsLock.sol           # Core contract orchestrating escrow state machine
│   └── Model.sol               # Shared structs, enums, events, and constants
└── test/                       # Forge-based Contract Tests
    ├── FundsLock.t.sol         # Happy/Success path unit tests
    ├── FundsLockReverts.t.sol  # Revert/Unhappy path unit tests
    └── TestHelper.t.sol        # Test suite base setups, mock utilities, and helpers
```

---

## 3. Domain & State Machine Model

Escrow agreements adhere to a strict state machine defined in `src/Model.sol` and managed in `src/FundsLock.sol`.

### Agreement State Flow (`AgreementStatus`)
1. **`CREATED`**: An agreement is registered by either the seller or the buyer specifying the counterparty and the `amount`.
   - If created by the seller, it is marked as `sellerAccepted = true`.
   - If created by the buyer, `sellerAccepted` remains `false` until the seller explicitly accepts.
2. **`SELLER_ACCEPTED`**: The seller accepts the terms of the agreement.
3. **`FUNDED`**: The buyer funds the exact agreement amount into the contract balance.
4. **`RELEASED`**: Funds are released to the seller after the buyer's approval or mediator resolution.
5. **`CANCELLED`**: The agreement is cancelled/refunded by mutual stakeholder agreement or mediator intervention.

### Escrow Core Struct (`EscrowAgreement`)
```solidity
struct EscrowAgreement {
    address seller;
    address payable buyer;
    uint256 amount;
    bool funded;
    bool sellerAccepted;
    bool sellerRequestedRelease;
    bool buyerApprovedRelease;
    bool released;
}
```

### Access Roles
- **Stakeholders:** The specified `seller` and `buyer`. Only these addresses can trigger standard agreement lifecycle actions (e.g., accept, fund, request release).
- **Mediator (`MEDIATOR_ROLE`):** Role assigned on contract construction (using OpenZeppelin `AccessControl`). Mediators resolve disputes and can override status transitions or trigger forced releases/refunds.

---

## 4. Development & Operational Runbook

Our `Makefile` orchestrates development tasks. You are expected to use these targets to validate any code modifications.

### Environment Preparation
Always initialize the project environment first. This creates a python virtual environment (`.venv`) for security analysis and pulls git submodules:
```bash
make init-project
```

### Main Verification Pipeline
Before submitting any changes, you must ensure they pass our build pipeline:
```bash
# Formats, lints, checks licenses, and runs tests
make build
```

Alternatively, run individual checks:
```bash
# Check code formatting (Forge)
make format-check

# Check Solidity style & quality guidelines
make lint

# Run all Unit Tests (Forge)
make test

# Generate a code coverage report for business logic contracts
make coverage
```

### Integration Testing with Kurtosis
Our integration tests spin up a local multi-node Ethereum testnet in a Kurtosis enclave, deploy the contracts, and execute transaction scenarios.
```bash
# Spins up Kurtosis and executes the script/IntegrationTest.s.sol scenario
make integration-test
```

### Security & Static Analysis
We use Slither to audit smart contracts against known vulnerabilities.
```bash
# Performs static analysis on src/ contracts
make static-analysis
```

---

## 5. Coding Standards & Conventions

AI agents must strictly adhere to these standards when writing or modifying code in this workspace:

### Strict File Header Rule
Every single `.sol` file under `src/`, `test/`, and `script/` **must** begin with the exact SPDX and Copyright headers. Failure to include these will cause `make copyright-check` to fail.
```solidity
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED
```

### Code Formatting
Ensure formatting adheres to Foundry defaults:
```bash
forge fmt
```

### Custom Error Patterns
Never use inline error strings or `require(condition, "error description")` statements. Define all error types inside `src/Errors.sol` and throw them using `revert`:
* *Incorrect:* `require(msg.value == amount, "Invalid amount");`
* *Correct:* `if (msg.value != amount) revert FundsLock_InvalidAmount(amount, msg.value);`

### Checks-Effects-Interactions (CEI)
To prevent reentrancy and state inconsistency, always structure functions using the CEI pattern:
1. **Checks:** Validate arguments, stakeholder identity, state conditions, and values.
2. **Effects:** Update contract state variables (e.g. `agreement.funded = true`, `agreement.released = true`).
3. **Interactions:** Trigger external calls, native ether transfers, or event emissions.

---

## 6. Testing & Validation Strategy

Our testing strategy mandates clear separation between happy path scenarios and unhappy path edge-cases.

1. **Unit Tests (Success Paths):** Implement successful transaction scenarios in `test/FundsLock.t.sol`. Use the `TestHelper` contract to create test wallets (`buyerWallet`, `sellerWallet`) and assign balances using Foundry cheatcodes.
2. **Unit Tests (Revert Paths):** Implement and assert all revert conditions inside `test/FundsLockReverts.t.sol` using `vm.expectRevert`.
3. **Integration Scenarios:** Write large, sequential multi-step integration scenarios in `script/IntegrationTest.s.sol`. Ensure your contract deployments in deployment scripts properly handle local mock networks vs. live public nodes.

---

## 7. Deployment Rules

### Local Kurtosis Deployment
Local deployments are orchestrated via:
```bash
make deploy-contract-kurtosis
```
This target spins up a Kurtosis Ethereum package, waits for block generation, and runs `script/deploy/DeployLocal.s.sol` to broadcast contract creation transactions.

### Public Deployment with AWS KMS
For public testnets and mainnets, we do **not** use raw private keys stored on-disk. Transactions are signed using AWS KMS:
1. Initialize/Import the key securely using `infra/secure-kms-import.sh`.
2. Define the target network RPC and AWS credentials:
   - `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_KMS_KEY_ID`, and `NETWORK_RPC_URL`.
3. Execute deployment:
   ```bash
   make deploy-contract-public
   ```
