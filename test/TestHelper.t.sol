// SPDX-License-Identifier: MPL-2.0fo
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {Test} from "forge-std/Test.sol";
import {FundsLock, EscrowAgreement, AgreementStatus, AgreementEvent} from "../src/FundsLock.sol";
import "../src/Errors.sol";
import {console} from "forge-std/console.sol";

struct Wallet {
    address addr;
    uint256 balance;
}

contract TestHelper is Test {
    function createWallet(string memory addr, uint256 balance) public returns (Wallet memory) {
        console.log("Creating wallet for ", addr, "with balance:", balance);
        return Wallet(makeAddr(addr), balance);
    }
}

