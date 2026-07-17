// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {Test} from "forge-std/Test.sol";
import {FundsLock, EscrowAgreement, AgreementStatus, AgreementEvent} from "../src/FundsLock.sol";
import "../src/Errors.sol";
import {console} from "forge-std/console.sol";
import {Wallet, TestHelper} from "./TestHelper.t.sol";

contract FundsLockTest is Test {
    FundsLock private fundsLock;
    TestHelper private helper;

    Wallet private buyerWallet;
    Wallet private sellerWallet;

    function setUp() public {
        fundsLock = new FundsLock();
        helper = new TestHelper();

        buyerWallet = helper.createWallet("buyer", 1 ether);
        sellerWallet = helper.createWallet("seller", 1 ether);

        vm.deal(buyerWallet.addr, buyerWallet.balance);
        vm.deal(sellerWallet.addr, sellerWallet.balance);
    }

    function test_CreateAgreementSuccess() public {
        uint256 testAmount = 0.5 ether;

        vm.expectEmit(true, true, false, true);
        emit AgreementEvent(sellerWallet.addr, buyerWallet.addr, testAmount, AgreementStatus.CREATED, block.timestamp);

        vm.startPrank(buyerWallet.addr);
        console.log("Creating agreement with seller:", sellerWallet.addr, "and buyer:", buyerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);
        vm.stopPrank();

        EscrowAgreement memory agreement = fundsLock.getAgreement(agreementId); // idCounter is incremented after returning

        assertEq(agreement.seller, sellerWallet.addr, "Seller address mismatch");
        assertEq(agreement.buyer, buyerWallet.addr, "Buyer address mismatch");
        assertEq(agreement.amount, testAmount, "Amount mismatch");
        assertEq(agreement.funded, false, "Agreement should not be funded initially");
        assertEq(agreement.sellerAccepted, false, "Seller should not have accepted initially");
        assertEq(agreement.sellerRequestedRelease, false, "Seller should not have requested release initially");
        assertEq(agreement.buyerApprovedRelease, false, "Buyer should not have approved release initially");
    }

    function test_SellerAcceptsAgreementAutomatically() public {
        uint256 testAmount = 0.5 ether;

        vm.expectEmit(true, true, false, true);
        emit AgreementEvent(sellerWallet.addr, buyerWallet.addr, testAmount, AgreementStatus.CREATED, block.timestamp);
        vm.expectEmit(true, true, false, true);
        emit AgreementEvent(
            sellerWallet.addr, buyerWallet.addr, testAmount, AgreementStatus.SELLER_ACCEPTED, block.timestamp
        );

        vm.startPrank(sellerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);
        vm.stopPrank();

        EscrowAgreement memory agreement = fundsLock.getAgreement(agreementId); // idCounter is incremented after returning

        assertEq(agreement.sellerAccepted, true, "Seller should have accepted automatically");
    }

    function test_SellerAcceptsAgreementSuccess() public {
        uint256 testAmount = 0.5 ether;

        vm.prank(buyerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);

        vm.expectEmit(true, true, false, true);
        emit AgreementEvent(
            sellerWallet.addr, buyerWallet.addr, testAmount, AgreementStatus.SELLER_ACCEPTED, block.timestamp
        );
        vm.prank(sellerWallet.addr);
        fundsLock.sellerAcceptAgreement(agreementId);
    }

    function test_BuyerFundsAgreementSuccess() public {
        uint256 testAmount = 0.5 ether;
        uint256 startingBalance = fundsLock.getBalance();

        vm.prank(buyerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);

        vm.expectEmit(true, true, false, true);
        emit AgreementEvent(sellerWallet.addr, buyerWallet.addr, testAmount, AgreementStatus.FUNDED, block.timestamp);

        vm.startPrank(buyerWallet.addr);
        fundsLock.fundAgreement{value: testAmount}(agreementId);
        vm.stopPrank();

        EscrowAgreement memory agreement = fundsLock.getAgreement(agreementId);
        assertEq(agreement.funded, true, "Agreement should be funded");
        assertEq(address(fundsLock).balance, startingBalance + testAmount, "Contract should have received the funds");
    }
}
