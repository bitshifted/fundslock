// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import "./Model.sol";
import "./Errors.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

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
        // addresses must be provided
        if (seller == address(0) || buyer == address(0)) {
            revert AddressNotSet();
        }
        // seller and buyer can not be the same address
        if (seller == buyer) {
            revert SameStakeholderAddress(seller);
        }
        // amount must be provided
        if (amount <= 0) {
            revert InsufficientAmount({minimum: 1, provided: amount});
        }
        // only seller or buyer can create agreement
        if (msg.sender != seller && msg.sender != buyer) {
            revert FundsLock_InvalidStakeholderAddress({provided: msg.sender, seller: seller, buyer: buyer});
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
        return idCounter - 1; // return the id of the newly created agreement (current id - 1)
    }

    function sellerAcceptAgreement(uint256 id) public {
        EscrowAgreement storage agreement = agreements[id];
        if (!agreementFound(agreement)) {
            revert FundsLock_AgreementNotFound(id);
        }
        if (msg.sender != agreement.seller) {
            revert FundsLock_InvalidStakeholderAddress(msg.sender, agreement.seller, agreement.buyer);
        }
        if (agreement.sellerAccepted) {
            revert FundsLock_AlreadyAccepted(id);
        }
        agreement.sellerAccepted = true;
        emit AgreementEvent(
            agreement.seller, agreement.buyer, agreement.amount, AgreementStatus.SELLER_ACCEPTED, block.timestamp
        );
    }

    function fundAgreement(uint256 id) public payable {
        EscrowAgreement storage agreement = agreements[id];

        // Checks
        if (!agreementFound(agreement)) {
            revert FundsLock_AgreementNotFound(id);
        }
        if (msg.sender != agreement.buyer) {
            revert FundsLock_InvalidStakeholderAddress(msg.sender, agreement.seller, agreement.buyer);
        }
        if (agreement.funded) {
            revert FundsLock_InvalidAgreementStatusTransition(AgreementStatus.FUNDED, AgreementStatus.FUNDED);
        }
        if (msg.value != agreement.amount) {
            revert FundsLock_InvalidAmount(agreement.amount, msg.value);
        }

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
        if (!agreementFound(agreement)) {
            revert FundsLock_AgreementNotFound(id);
        }
        // only seller or buyer can release funds
        if (msg.sender != agreement.seller && msg.sender != agreement.buyer) {
            revert FundsLock_InvalidStakeholderAddress(msg.sender, agreement.seller, agreement.buyer);
        }
        if (!agreement.funded) {
            revert FundsLock_InvalidAgreementStatusTransition(AgreementStatus.FUNDED, AgreementStatus.RELEASED);
        }
        if (!agreement.sellerAccepted) {
            revert FundsLock_InvalidAgreementStatusTransition(AgreementStatus.SELLER_ACCEPTED, AgreementStatus.RELEASED);
        }

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
