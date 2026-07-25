// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {Test} from "forge-std/Test.sol";
import {Vm} from "forge-std/Vm.sol";
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

    function assertAgreementCreatedLog(
        Vm.Log memory log,
        address expectedEmitter,
        uint256 expectedId,
        address expectedSeller,
        address expectedBuyer,
        uint256 expectedAmount
    ) internal view {
        assertEq(log.emitter, expectedEmitter, "Unexpected emitter");
        assertEq(uint256(log.topics[1]), expectedId, "Unexpected agreement id");
        assertEq(address(uint160(uint256(log.topics[2]))), expectedSeller, "Unexpected seller");
        assertEq(address(uint160(uint256(log.topics[3]))), expectedBuyer, "Unexpected buyer");

        (uint256 amount, uint256 timestamp) = abi.decode(log.data, (uint256, uint256));
        assertEq(amount, expectedAmount, "Unexpected amount");
        assertEq(timestamp, block.timestamp, "Unexpected timestamp");
    }

    function assertAgreementEventLog(
        Vm.Log memory log,
        address expectedEmitter,
        uint256 expectedId,
        address expectedSeller,
        address expectedBuyer,
        AgreementStatus expectedStatus
    ) internal view {
        assertEq(log.emitter, expectedEmitter, "Unexpected emitter");
        assertEq(address(uint160(uint256(log.topics[1]))), expectedSeller, "Unexpected seller");
        assertEq(address(uint160(uint256(log.topics[2]))), expectedBuyer, "Unexpected buyer");
        assertEq(uint256(log.topics[3]), expectedId, "Unexpected agreement id");

        (uint8 statusValue, uint256 timestamp) = abi.decode(log.data, (uint8, uint256));
        assertEq(uint8(expectedStatus), statusValue, "Unexpected status");
        assertEq(timestamp, block.timestamp, "Unexpected timestamp");
    }

    function test_CreateAgreementSuccess() public {
        uint256 testAmount = 0.5 ether;

        vm.recordLogs();
        vm.startPrank(buyerWallet.addr);
        console.log("Creating agreement with seller:", sellerWallet.addr, "and buyer:", buyerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);
        vm.stopPrank();

        Vm.Log[] memory entries = vm.getRecordedLogs();
        assertEq(entries.length, 1, "Expected one event");
        assertAgreementCreatedLog(
            entries[0], address(fundsLock), agreementId, sellerWallet.addr, buyerWallet.addr, testAmount
        );

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

        vm.recordLogs();
        vm.startPrank(sellerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);
        vm.stopPrank();

        Vm.Log[] memory entries = vm.getRecordedLogs();
        assertEq(entries.length, 2, "Expected two events");
        assertAgreementCreatedLog(
            entries[0], address(fundsLock), agreementId, sellerWallet.addr, buyerWallet.addr, testAmount
        );
        assertAgreementEventLog(
            entries[1],
            address(fundsLock),
            agreementId,
            sellerWallet.addr,
            buyerWallet.addr,
            AgreementStatus.SELLER_ACCEPTED
        );

        EscrowAgreement memory agreement = fundsLock.getAgreement(agreementId); // idCounter is incremented after returning

        assertEq(agreement.sellerAccepted, true, "Seller should have accepted automatically");
    }

    function test_SellerAcceptsAgreementSuccess() public {
        uint256 testAmount = 0.5 ether;

        vm.prank(buyerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);

        vm.recordLogs();
        vm.prank(sellerWallet.addr);
        fundsLock.sellerAcceptAgreement(agreementId);

        Vm.Log[] memory entries = vm.getRecordedLogs();
        assertEq(entries.length, 1, "Expected one event");
        assertAgreementEventLog(
            entries[0],
            address(fundsLock),
            agreementId,
            sellerWallet.addr,
            buyerWallet.addr,
            AgreementStatus.SELLER_ACCEPTED
        );
    }

    function test_BuyerFundsAgreementSuccess() public {
        uint256 testAmount = 0.5 ether;
        uint256 startingBalance = fundsLock.getBalance();

        vm.prank(buyerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);

        vm.recordLogs();
        vm.startPrank(buyerWallet.addr);
        fundsLock.fundAgreement{value: testAmount}(agreementId);
        vm.stopPrank();

        Vm.Log[] memory entries = vm.getRecordedLogs();
        assertEq(entries.length, 1, "Expected one event");
        assertAgreementEventLog(
            entries[0], address(fundsLock), agreementId, sellerWallet.addr, buyerWallet.addr, AgreementStatus.FUNDED
        );

        EscrowAgreement memory agreement = fundsLock.getAgreement(agreementId);
        assertEq(agreement.funded, true, "Agreement should be funded");
        assertEq(address(fundsLock).balance, startingBalance + testAmount, "Contract should have received the funds");
    }

    function test_RequestReleaseSuccess() public {
        uint256 testAmount = 0.5 ether;

        vm.prank(buyerWallet.addr);
        uint256 agreementId = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);

        vm.prank(sellerWallet.addr);
        fundsLock.sellerAcceptAgreement(agreementId);

        vm.recordLogs();
        vm.prank(sellerWallet.addr);
        fundsLock.requestRelease(agreementId);

        Vm.Log[] memory entries = vm.getRecordedLogs();
        assertEq(entries.length, 1, "Expected one event");
        assertAgreementEventLog(
            entries[0],
            address(fundsLock),
            agreementId,
            sellerWallet.addr,
            buyerWallet.addr,
            AgreementStatus.SELLER_REQUESTED_RELEASE
        );

        EscrowAgreement memory agreement = fundsLock.getAgreement(agreementId);
        assertEq(agreement.sellerRequestedRelease, true, "Seller should have requested release");
    }

    function test_ReleaseFundsSuccess() public {
        uint256 testAmount = 0.5 ether;

        // 1. Create agreement
        vm.prank(buyerWallet.addr);
        uint256 id = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);

        // 2. Seller accepts
        vm.prank(sellerWallet.addr);
        fundsLock.sellerAcceptAgreement(id);

        // 3. Buyer funds
        vm.prank(buyerWallet.addr);
        fundsLock.fundAgreement{value: testAmount}(id);

        // Check pre-release balances
        uint256 sellerStartingBalance = sellerWallet.addr.balance;

        // 4. Seller requests release
        vm.prank(sellerWallet.addr);
        fundsLock.releaseFunds(id);

        // 5. Buyer approves release
        vm.recordLogs();
        vm.prank(buyerWallet.addr);
        fundsLock.releaseFunds(id);

        Vm.Log[] memory entries = vm.getRecordedLogs();
        assertEq(entries.length, 1, "Expected one event");
        assertAgreementEventLog(
            entries[0], address(fundsLock), id, sellerWallet.addr, buyerWallet.addr, AgreementStatus.RELEASED
        );

        // 6. Verify results
        EscrowAgreement memory agreement = fundsLock.getAgreement(id);
        assertEq(agreement.released, true, "Agreement should be marked as released");
        assertEq(sellerWallet.addr.balance, sellerStartingBalance + testAmount, "Seller should have received the funds");
    }
}
