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
        string memory networkName = vm.envString("NETWORK_NAME");
        console.log("Deploying to network: ", networkName);
        vm.startBroadcast();
        fundsLock = new FundsLock();
        vm.stopBroadcast();
        console.log("Contract deployed to address: ", vm.toString(address(fundsLock)));
        console.log("Deployed at block number: ", block.number);

        // copy ABI to Graph protocol config
        string memory abiContent = vm.readFile("out/FundsLock.sol/FundsLock.json");
        vm.writeFile("../graph/abis/FundsLock.json", abiContent);
        console.log("ABI copied to Graph protocol config");
        
        // update values in Graph protocol networks config
        string[] memory cmds = new string[](3);
        cmds[0] = "bash";
        cmds[1] = "-c";
        cmds[2] = string.concat(
            "cd ../graph && ./update_networks.sh ",
            networkName,
            " ",
            vm.toString(address(fundsLock)),
            " ",
            vm.toString(block.number)
        );
        vm.ffi(cmds);
        console.log("Graph protocol networks config updated for network: ", networkName);
    }
}
