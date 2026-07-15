// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

pragma solidity ^0.8.35;

import {Script} from "forge-std/Script.sol";
import {FundsLock} from "../../src/FundsLock.sol";
import {stdJson} from "forge-std/StdJson.sol";
import {console} from "forge-std/console.sol";

/**
 * Deploys contract to public tesnets/mainnets, like Sepolia or Ethereum mainnet. It is intended to run with KMS as deployment
 * key provider.
 *
 * @title DeployPublic
 * @dev A script to deploy the  contract to public networks using KMS for signing
 */
contract DeployPublic is Script {
    FundsLock internal fundsLock;

    function run() public {
        vm.startBroadcast();
        fundsLock = new FundsLock();
        vm.stopBroadcast();
        console.log("Contract deployed to address: ", vm.toString(address(fundsLock)));
    }
}
