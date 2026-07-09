// SPDX-Licrnse-Identifier: MIT

pragma solidity ^0.8.35;

import {Script} from "forge-std/Script.sol";
// import {BasicEscrow, EscrowAgreement} from "../../../src/BasicEscrow.sol";
import {console} from "forge-std/console.sol";
import {Walets} from "./common/Walets.s.sol";
import {Test} from "forge-std/Test.sol";

contract IntegrationTest is Script, Test {
    // BasicEscrow private basicEscrow;
    Walets private wallets;

    function setUp() public {
        wallets = new Walets("./config/test-accounts.json");
        wallets.run();
        string memory addrString = vm.readFile("./testout/address.txt");
        console.log("BasicEscrow contract address loaded from testout/address.txt: ", addrString);
        // basicEscrow = BasicEscrow(payable(vm.parseAddress(addrString)));
    }

    function run() public {
        // createAndFundAgreementSuccess(0.5 ether);
    }

    // function createAndFundAgreementSuccess(uint256 amount) public {
    //     vm.startBroadcast(wallets.forRole(ParticipantRole.Buyer).privateKey);
    //     uint256 agreementId = basicEscrow.createAgreement(
    //         wallets.forRole(ParticipantRole.Seller).addr, payable(wallets.forRole(ParticipantRole.Buyer).addr), amount
    //     );
    //     console.log("Agreement created with ID: ", agreementId);
    //     vm.stopBroadcast();
    //     EscrowAgreement memory agreement = basicEscrow.getAgreement(agreementId);
    //     assertEq(agreement.seller, wallets.forRole(ParticipantRole.Seller).addr, "Seller address mismatch");
    //     assertEq(agreement.buyer, wallets.forRole(ParticipantRole.Buyer).addr, "Buyer address mismatch");
    //     assertEq(agreement.amount, amount, "Amount mismatch");
    //     assertEq(agreement.funded, false, "Agreement should not be funded initially");

    //     vm.startBroadcast(wallets.forRole(ParticipantRole.Seller).privateKey);
    //     basicEscrow.sellerAcceptAgreement(agreementId);
    //     vm.stopBroadcast();
    //     console.log("Seller accepted the agreement with ID: ", agreementId);
    // }
}
