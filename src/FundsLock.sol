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
        // only sender ot buyer can create agreement
        require(msg.sender == seller || msg.sender == buyer, "Only the seller or buyer can create an agreement");
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
}
