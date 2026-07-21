// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

enum AgreementStatus {
    CREATED,
    FUNDED,
    SELLER_ACCEPTED,
    RELEASED,
    CANCELLED
}

struct EscrowAgreement {
    address seller;
    bool funded;
    bool sellerAccepted;
    bool sellerRequestedRelease;
    bool buyerApprovedRelease;
    bool released;
    address payable buyer;
    uint256 amount;
}

event AgreementEvent(
    address indexed seller, address indexed buyer, uint256 amount, AgreementStatus status, uint256 timestamp
);
