// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {Test} from "forge-std/Test.sol";
import {FundsLock, EscrowAgreement, AgreementStatus, AgreementEvent} from "../src/FundsLock.sol";
import "../src/Errors.sol";
import {console} from "forge-std/console.sol";

contract FundsLockTest is Test {
    FundsLock public fundsLock;

    struct Wallet {
        address addr;
        uint256 balance;
    }

    Wallet public buyerWallet;
    Wallet public sellerWallet;

    function setUp() public {
        fundsLock = new FundsLock();

        buyerWallet = Wallet({addr: makeAddr("buyer"), balance: 1 ether});
        sellerWallet = Wallet({addr: makeAddr("seller"), balance: 1 ether});

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

    function test_CreateAgreementFailsWhenSenderMismatch() public {
        uint256 testAmount = 0.5 ether;
        Wallet memory scamWallet = Wallet({addr: makeAddr("scam"), balance: 1 ether});
        vm.deal(scamWallet.addr, scamWallet.balance);

        vm.expectRevert(
            abi.encodeWithSelector(
                InvalidStakeholderAddress.selector, scamWallet.addr, sellerWallet.addr, buyerWallet.addr
            )
        );
        vm.prank(scamWallet.addr);
        fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);
    }
}
