// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {Script} from "forge-std/Script.sol";
import {FundsLock, EscrowAgreement} from "../src/FundsLock.sol";
import {console} from "forge-std/console.sol";
import {Wallets} from "./common/Wallets.s.sol";
import {Test} from "forge-std/Test.sol";

contract IntegrationTest is Script, Test {
    FundsLock private fundsLock;
    Wallets private wallets;

    function setUp() public {
        wallets = new Wallets();
        wallets.run("./config/test-accounts.json");
        string memory addrString = vm.readFile("./testout/address.txt");
        console.log("FundsLock contract address loaded from testout/address.txt: ", addrString);
        fundsLock = FundsLock(payable(vm.parseAddress(addrString)));
    }

    function run() public {
        createAndFundAgreementSuccess(0.5 ether);
    }

    function createAndFundAgreementSuccess(uint256 amount) public {
        vm.startBroadcast(wallets.forRole("buyer").privateKey);
        uint256 agreementId =
            fundsLock.createAgreement(wallets.forRole("seller").addr, payable(wallets.forRole("buyer").addr), amount);
        console.log("Agreement created with ID: ", agreementId);
        vm.stopBroadcast();
        EscrowAgreement memory agreement = fundsLock.getAgreement(agreementId);
        assertEq(agreement.seller, wallets.forRole("seller").addr, "Seller address mismatch");
        assertEq(agreement.buyer, wallets.forRole("buyer").addr, "Buyer address mismatch");
        assertEq(agreement.amount, amount, "Amount mismatch");
        assertEq(agreement.funded, false, "Agreement should not be funded initially");

        vm.startBroadcast(wallets.forRole("seller").privateKey);
        fundsLock.sellerAcceptAgreement(agreementId);
        vm.stopBroadcast();
        console.log("Seller accepted the agreement with ID: ", agreementId);
    }
}
