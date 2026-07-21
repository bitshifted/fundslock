// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import "./Model.sol";

error AddressNotSet();

error InsufficientAmount(uint256 minimum, uint256 provided);

error FundsLock_InvalidStakeholderAddress(address provided, address seller, address buyer);

error SameStakeholderAddress(address provided);

error FundsLock_AgreementNotFound(uint256 id);

error FundsLock_AlreadyAccepted(uint256 id);

error FundsLock_InvalidAgreementStatusTransition(AgreementStatus current, AgreementStatus next);

error FundsLock_InvalidAmount(uint256 required, uint256 provided);

