// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.35;

import {Script} from "forge-std/Script.sol";
import {FundsLock} from "../../src/FundsLock.sol";
import {stdJson} from "forge-std/StdJson.sol";
import {console} from "forge-std/console.sol";
import {Wallets} from "../common/Wallets.s.sol";

/**
 * Deploys contract to local network, like Anvil or Kurtosis. It takes a predefined set of wallets from a JSON file
 * located at ./config/test-accounts.json and deploys the  contract. The deployed contract address is saved to a file .testout/address.txt
 * for later use in integration tests.
 * The main purpose of this script is to facilitate local testing and impprove developer experience, not secure deployment.
 *
 * @title DeployLocal
 * @dev A script to deploy the  contract locally for testing.
 */
contract DeployLocal is Script {
    FundsLock internal fundsLock;
    Wallets private wallets;

    function setUp() public {
        wallets = new Wallets();
        wallets.run("./config/test-accounts.json");
    }

    function run() public {
        vm.startBroadcast(wallets.forRole("deployer").privateKey);
        fundsLock = new FundsLock();
        vm.stopBroadcast();
        vm.writeFile("./testout/address.txt", vm.toString(address(fundsLock)));
    }

    function getFundsLock() public view returns (FundsLock) {
        return fundsLock;
    }

    function getWallets() public view returns (Wallets) {
        return wallets;
    }
}
