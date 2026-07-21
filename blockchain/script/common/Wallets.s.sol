// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Bitshift ED

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
 * Loads addresses and private keys from a JSON configuration file for testing and deployment purposes.
 * @title Wallets
 * @dev holds addresses and keys loaded from configuration JSON file
 */
contract Wallets is Script {
    mapping(string => Wallet) internal wallets;

    /**
     * @dev Loads wallets from a JSON configuration file. The JSON file should have the following structure:
     * {
     *   "role1": {
     *     "address": "0x...",
     *     "private_key": "..."
     *   },
     *   "role2": {
     *     "address": "0x...",
     *     "private_key": "..."
     *   }
     * }
     * Each top-level key represents a role, and the corresponding value is an object containing the address and private key for that role..
     * @param configPath path where configuration file is located
     */
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

    /**
     * @dev Returns the wallet for a given role.
     * @param role role for which to get the wallet
     */
    function forRole(string memory role) public view returns (Wallet memory) {
        return wallets[role];
    }
}
