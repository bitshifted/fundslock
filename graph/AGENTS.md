# Agent Guidelines & Runbook (AGENTS.md)

Welcome, AI Agent! This file serves as your comprehensive instruction manual and onboarding guide for working in the **graph** directory. It contains our architectural layout, operational commands, coding guidelines, and deployment procedures to ensure your contributions align perfectly with our standards.

---

## 1. Repository Overview

This directory contains a [The Graph](https://thegraph.com/) protocol subgraph that indexes events from the FundsLock smart contract. The subgraph enables efficient querying of escrow agreement events and state changes through GraphQL.

### Key Technology Stack
- **Subgraph Framework:** The Graph Protocol (`@graphprotocol/graph-cli`, `@graphprotocol/graph-ts`)
- **Language:** AssemblyScript (via `wasm/assemblyscript` mapping)
- **Schema Definition:** GraphQL (`schema.graphql`)
- **Testing:** Matchstick (`matchstick-as`)
- **Local Development:** Docker Compose (IPFS + Graph Node + Indexer Node)

---

## 2. Directory Architecture

```
/graph/
├── abis/FundsLock.json             # ABI extracted from deployed smart contract
├── src/funds-lock.ts               # Event handlers (mapping logic in AssemblyScript)
├── schema.graphql                  # GraphQL entity definitions
├── subgraph.yaml                   # Subgraph manifest (data sources, event mappings)
├── networks.json                   # Contract addresses and start blocks per network
├── tests/
│   ├── funds-lock.test.ts          # Matchstick unit tests
│   └── funds-lock-utils.ts         # Test mock event factories
├── docker-compose.yml              # Local development environment
├── package.json                    # Dependencies and scripts
├── tsconfig.json                   # TypeScript/AssemblyScript config
├── update_networks.sh              # Script to refresh networks.json from deployments
└── AGENTS.md                       # This file
```

---

## 3. Domain Model & Entities

### AgreementLog Entity

The subgraph indexes a single entity type: `AgreementLog`. Each logged event corresponds to an escrow agreement lifecycle event emitted by the `FundsLock` smart contract.

```graphql
type AgreementLog @entity(immutable: true) {
  id: Bytes!

  agreement_id: BigInt!    # uint256 - unique agreement identifier
  seller: Bytes!           # address - seller's Ethereum address
  buyer: Bytes!            # address - buyer's Ethereum address
  amount: BigInt!          # uint256 - funded amount (set only on AgreementCreated)
  status: Int!             # uint8 - current agreement status
  timestamp: BigInt!       # uint256 - event timestamp

  blockNumber: BigInt!     # block number where event was emitted
  blockTimestamp: BigInt!  # block timestamp
  transactionHash: Bytes!  # transaction hash containing the event
}
```

The entity is **immutable** — once created, it cannot be modified. Each event produces a new `AgreementLog` instance identified by `transactionHash.logIndex`.

### Tracked Events

| Event | Handler | Notes |
|-------|---------|-------|
| `AgreementCreated(indexed uint256, indexed address, indexed address, uint256, uint256)` | `handleAgreementCreated` | Initializes agreement; sets `amount`, initial `status = 0` |
| `AgreementEvent(indexed address, indexed address, indexed uint256, uint8, uint256)` | `handleAgreementEvent` | Tracks status transitions; `amount = 0` (not applicable here) |

---

## 4. Development & Operational Commands

All commands should be run from the `/graph` directory.

### Prerequisites
Install dependencies:
```bash
yarn install
```

### Code Generation
After modifying `schema.graphql` or `subgraph.yaml`, regenerate TypeScript bindings:
```bash
yarn codegen
```
Or using the CLI directly:
```bash
npx graph codegen
```

### Build
Compile the subgraph into WASM:
```bash
NETWORK=<network_name> yarn build
```
The `NETWORK` variable specifies which network configuration to use from `networks.json`. For example:
```bash
NETWORK=sepolia yarn build
```

### Local Development & Testing

Spin up local Graph Node stack:
```bash
docker-compose up -d
```

Create and deploy locally:
```bash
yarn create-local        # Register subgraph locally
yarn deploy-local         # Build and deploy to local node
yarn remove-local         # Remove local subgraph
```

Run tests:
```bash
yarn test
```

### Deployment

Deployment is automated via GitHub Actions. Before merging:

1. Ensure `networks.json` reflects the correct contract address and start block for the target network.
2. Make sure all tests pass: `yarn test`
3. Push to a feature branch — CI will build and test the subgraph.
4. Merge the PR — GitHub Actions will deploy the subgraph automatically.

Deployment uses the `$NETWORK` and `${GRAPH_SLUG}` environment variables set in the CI workflow.

### Publishing

Subgraphs must be published through [Subgraph Studio](https://thegraph.com/studio/). This step requires GRT tokens in the wallet and cannot be done via CI/CD.

---

## 5. Coding Standards & Conventions

### AssemblyScript Rules

- Use imports from `@graphprotocol/graph-ts` for all primitive types (`BigInt`, `Bytes`, `Address`).
- Entity IDs must be derived deterministically. Current pattern: `event.transaction.hash.concatI32(event.logIndex.toI32())`.
- All fields in the `AgreementLog` entity must be set before calling `.save()`.
- Follow existing naming conventions: handler functions are named `handle<EventName>` matching the event name from `subgraph.yaml`.

### Schema Changes

When adding new entities or fields to `schema.graphql`:
1. Update `schema.graphql` first.
2. Run `yarn codegen` to regenerate TypeScript types.
3. Update `subgraph.yaml` if adding new data sources or event handlers.
4. Run `yarn build` to verify compilation.

### Networks Configuration

The `networks.json` file maps network names to contract deployment data. When the smart contract is redeployed to a new network, run:
```bash
./update_networks.sh
```
This script auto-populates `networks.json` from deployment artifacts.

### Test Conventions

- Tests live in `tests/` alongside their source files.
- Use `matchstick-as`'s `newMockEvent()` to create mock events.
- Create helper factory functions in `tests/<source>-utils.ts` for reusable mock event builders.
- Assertions use standard `assert.*` functions from `matchstick-as`.

---

## 6. Troubleshooting

### Subgraph Purged (Not Indexed)

If queries return errors like `deployment <slug>/<version>/latest does not exist`, the subgraph may have been purged due to inactivity in the dev environment. Fix by making a small commit and pushing to trigger a fresh deployment.

### Build Failures

- Ensure `NETWORK` is set correctly and exists in `networks.json`.
- After changing `schema.graphql`, always run `yarn codegen` before building.
- Verify that event signatures in `subgraph.yaml` match the ABI exactly (parameter types and ordering).

### Address Mismatch

Contract addresses and start blocks in `networks.json` must match the actual deployed contract. Always regenerate this file after deployments rather than editing manually.
