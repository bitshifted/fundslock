# Agent Guidelines & Runbook (AGENTS.md)

his file serves as your comprehensive instruction manual and onboarding guide for working in the **fundslock** repository. It contains our architectural layout, domain model, operational commands, coding guidelines, and testing requirements.

## Project overview

**FundsLock** is decentralized escrow system implemented as an EVM smart contract, with off-chain backend and frontend component. It enables buyers and sellers to establish secure agreements, fund contracts, request/approve releases, and engage in mediator-driven dispute resolutions.

## Repository structure

The repository consists of the following component.

* `blockchain` - Solidity/Foundry project implementing smart contract for escrow
* `backend` - Go-based REST API server

Each component has it's own `AGENTS,md` file with instruction specific to the component.
