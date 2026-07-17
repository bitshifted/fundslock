// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {Test} from "forge-std/Test.sol";
import {console} from "forge-std/console.sol";
import {FundsLock} from "../src/FundsLock.sol";
import {Wallet, TestHelper} from "./TestHelper.t.sol";
import "../src/Model.sol";
import "../src/Errors.sol";

/**
 * @title FundsLockReverts
 * @notice Contains tests for invalid conditions
 */
contract FundsLockReverts is Test {
    FundsLock private fundsLock;
    TestHelper private helper;

    Wallet private buyerWallet;
    Wallet private sellerWallet;

    function setUp() public {
        fundsLock = new FundsLock();
        helper = new TestHelper();

        buyerWallet = helper.createWallet("buyer", 10 ether);
        sellerWallet = helper.createWallet("seller", 10 ether);

        vm.deal(buyerWallet.addr, buyerWallet.balance);
        vm.deal(sellerWallet.addr, sellerWallet.balance);
    }

    // createAgreement() tests

    function test_StakeholdersCanNotBeEmpty() public {
        vm.expectRevert(abi.encodeWithSelector(AddressNotSet.selector));
        vm.prank(sellerWallet.addr);
        fundsLock.createAgreement(address(0), payable(buyerWallet.addr), 1 ether);
    }

    function test_BuyerSellerCanNotBeSame() public {
        vm.expectRevert(abi.encodeWithSelector(SameStakeholderAddress.selector, sellerWallet.addr));
        vm.prank(sellerWallet.addr);
        fundsLock.createAgreement(sellerWallet.addr, payable(sellerWallet.addr), 1 ether);
    }

    function test_AmountMustBePositive() public {
        vm.expectRevert(abi.encodeWithSelector(InsufficientAmount.selector, 1, 0));
        vm.prank(buyerWallet.addr);
        fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), 0);
    }

    function test_CreateAgreementFailsWhenSenderMismatch() public {
        uint256 testAmount = 0.5 ether;
        Wallet memory scamWallet = helper.createWallet("scam", 1 ether);
        vm.deal(scamWallet.addr, scamWallet.balance);

        vm.expectRevert(
            abi.encodeWithSelector(
                FundsLock_InvalidStakeholderAddress.selector, scamWallet.addr, sellerWallet.addr, buyerWallet.addr
            )
        );
        vm.prank(scamWallet.addr);
        fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), testAmount);
    }

    // sellerAccepted() tests
    function test_SellerAcceptFailsWhenAgreementNotFound() public {
        vm.expectRevert(abi.encodeWithSelector(FundsLock_AgreementNotFound.selector, 12345));
        vm.prank(sellerWallet.addr);
        fundsLock.sellerAcceptAgreement(12345);
    }

    function test_OnlySellerCanAccept() public {
        Wallet memory scamWallet = helper.createWallet("scam", 1 ether);
        vm.deal(scamWallet.addr, scamWallet.balance);
        vm.prank(buyerWallet.addr);
        uint256 id = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), 1 ether);

        vm.expectRevert(
            abi.encodeWithSelector(
                FundsLock_InvalidStakeholderAddress.selector, scamWallet.addr, sellerWallet.addr, buyerWallet.addr
            )
        );
        vm.prank(scamWallet.addr);
        fundsLock.sellerAcceptAgreement(id);
    }

    function test_CanNotAcceptAlreadyAcceptedAgreement() public {
        vm.prank(buyerWallet.addr);
        uint256 id = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), 1 ether);

        vm.prank(sellerWallet.addr);
        fundsLock.sellerAcceptAgreement(id);

        vm.expectRevert(abi.encodeWithSelector(FundsLock_AlreadyAccepted.selector, id));
        vm.prank(sellerWallet.addr);
        fundsLock.sellerAcceptAgreement(id);
    }

    function test_FundAgreementFailsWhenAlreadyFunded() public {
        vm.prank(buyerWallet.addr);
        uint256 id = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), 1 ether);

        vm.prank(buyerWallet.addr);
        fundsLock.fundAgreement{value: 1 ether}(id);

        vm.expectRevert(
            abi.encodeWithSelector(
                FundsLock_InvalidAgreementStatusTransition.selector, AgreementStatus.FUNDED, AgreementStatus.FUNDED
            )
        );
        vm.prank(buyerWallet.addr);
        fundsLock.fundAgreement{value: 1 ether}(id);
    }

    function test_FundAgreementFailsWhenAmountIncorrect() public {
        uint256 amount = 1 ether;
        vm.prank(buyerWallet.addr);
        uint256 id = fundsLock.createAgreement(sellerWallet.addr, payable(buyerWallet.addr), amount);

        vm.expectRevert(abi.encodeWithSelector(FundsLock_InvalidAmount.selector, amount, amount / 2));
        vm.prank(buyerWallet.addr);
        fundsLock.fundAgreement{value: amount / 2}(id);
    }
}
