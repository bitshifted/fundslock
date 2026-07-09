// SPDX-License-Identifier: MIT

pragma solidity ^0.8.35;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";

/**
 * @title Wallet representation for testing and deployment
 * @dev represents a wallet with basic info like address and private key
 */
struct Wallet {
    address addr;
    uint256 privateKey;
}

/**
 * @title Walets
 * @dev holds addresses and keys loaded from configutation JSON file
 */
contract Walets is Script {
    mapping(string => Wallet) internal wallets;

    function run(string memory configPath) public {
        // load the deployer private key from the accounts.json file
        string memory json = vm.readFile(configPath);
        // list root-level keys, which represent the roles of the wallets
        string[] memory keys = vm.parseJsonKeys(json, "$");
        for (uint8 i = 0; i < keys.length; i++) {
            string memory key = keys[i];
            console.log("Loaded wallet role: ", key);
            wallets[key] = Wallet({
                addr: vm.parseJsonAddress(json, string.concat(".", key, ".address")),
                privateKey: vm.parseJsonUint(json, string.concat(".", key, ".private_key"))
            });
        }
        console.log("Wallets loaded from test-accounts.json:");
    }

    function forRole(string memory role) public view returns (Wallet memory) {
        return wallets[role];
    }
}
