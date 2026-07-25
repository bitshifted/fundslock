// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

enum AgreementStatus {
    CREATED,
    FUNDED,
    SELLER_ACCEPTED,
    SELLER_REQUESTED_RELEASE,
    BUYER_APPROVED_RELEASE,
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

event AgreementCreated(
    uint256 indexed id, address indexed seller, address indexed buyer, uint256 amount, uint256 timestamp
);

event AgreementEvent(
    address indexed seller, address indexed buyer, uint256 indexed id, AgreementStatus status, uint256 timestamp
);
