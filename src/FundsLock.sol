// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import "./Model.sol";
import "./Errors.sol";

contract FundsLock is AccessControl {
    bytes32 public constant MEDIATOR_ROLE = keccak256("MEDIATOR_ROLE");
    uint256 private idCounter;

    mapping(uint256 => EscrowAgreement) private agreements;

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MEDIATOR_ROLE, msg.sender);
        idCounter = 100;
    }

    function createAgreement(address seller, address payable buyer, uint256 amount) public returns (uint256) {
        if (seller == address(0) || buyer == address(0)) {
            revert AddressNotSet();
        }
        // amount must be provided
        if (amount <= 0) {
            revert InsufficientAmount({minimum: 1, provided: amount});
        }
        // only seller or buyer can create agreement
        if (msg.sender != seller && msg.sender != buyer) {
            revert InvalidStakeholderAddress({provided: msg.sender, seller: seller, buyer: buyer});
        }
        // if seller is the sender, it implies he already accepted the agreement
        bool sellerAccepted = msg.sender == seller ? true : false;
        EscrowAgreement memory agreement = EscrowAgreement({
            seller: seller,
            buyer: buyer,
            amount: amount,
            funded: false,
            sellerAccepted: sellerAccepted,
            sellerRequestedRelease: false,
            buyerApprovedRelease: false,
            released: false
        });
        agreements[idCounter++] = agreement;
        emit AgreementEvent(seller, buyer, amount, AgreementStatus.CREATED, block.timestamp);
        if (sellerAccepted) {
            emit AgreementEvent(seller, buyer, amount, AgreementStatus.SELLER_ACCEPTED, block.timestamp);
        }
        emit AgreementEvent(seller, buyer, amount, AgreementStatus.CREATED, block.timestamp);
        return idCounter - 1; // return the id of the newly created agreement (current id - 1)
    }

    function sellerAcceptAgreement(uint256 id) public {
        EscrowAgreement storage agreement = agreements[id];
        require(agreementFound(agreement), "Agreement not found");
        require(agreement.seller != address(0), "Agreement not found");
        require(msg.sender == agreement.seller, "Only the seller can accept the agreement");
        require(!agreement.sellerAccepted, "Agreement already accepted");
        agreement.sellerAccepted = true;
        emit AgreementEvent(
            agreement.seller, agreement.buyer, agreement.amount, AgreementStatus.SELLER_ACCEPTED, block.timestamp
        );
    }

    function fundAgreement(uint256 id) public payable {
        EscrowAgreement storage agreement = agreements[id];

        // Checks
        require(agreementFound(agreement), "Agreement not found");
        require(msg.sender == agreement.buyer, "Only the buyer can fund the agreement");
        require(!agreement.funded, "Agreement already funded");
        require(msg.value == agreement.amount, "Incorrect funding amount");

        // Effects: Funds are automatically transferred to the contract balance because the function is `payable`.
        agreement.funded = true;
        emit AgreementEvent(
            agreement.seller, agreement.buyer, agreement.amount, AgreementStatus.FUNDED, block.timestamp
        );
    }

    function agreementFound(EscrowAgreement storage agreement) private view returns (bool) {
        return agreement.seller != address(0) && agreement.buyer != address(0) && agreement.amount > 0;
    }

    function releaseFunds(uint256 id) public payable {
        EscrowAgreement storage agreement = agreements[id];
        require(agreementFound(agreement), "Agreement not found");
        // only seller or buyer can release funds
        require(
            msg.sender == agreement.seller || msg.sender == agreement.buyer,
            "Only the seller or buyer can release funds"
        );
        require(agreement.funded, "Agreement has not been funded");
        require(agreement.sellerAccepted, "Seller has not accepted the agreement");

        if (msg.sender == agreement.seller && !agreement.sellerRequestedRelease) {
            agreement.sellerRequestedRelease = true;
        }
        if (msg.sender == agreement.buyer && !agreement.buyerApprovedRelease) {
            agreement.buyerApprovedRelease = true;
        }
        if (agreement.sellerRequestedRelease && agreement.buyerApprovedRelease) {
            require(!agreement.released, "Funds already released");
            agreement.released = true;
            emit AgreementEvent(
                agreement.seller, agreement.buyer, agreement.amount, AgreementStatus.RELEASED, block.timestamp
            );
            payable(agreement.seller).transfer(agreement.amount);
            // (bool success,) = agreement.seller.call{value: agreement.amount}("");
            // require(success, "Failed to send Ether");
        }
    }

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }

    function getAgreement(uint256 id) public view returns (EscrowAgreement memory) {
        return agreements[id];
    }
}
