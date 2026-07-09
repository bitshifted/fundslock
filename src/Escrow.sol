// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";

contract Escrow is AccessControl {
    bytes32 public constant ESCROW_ROLE = keccak256("ESCROW_ROLE");

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ESCROW_ROLE, msg.sender);
    }
}
