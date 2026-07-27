import { newMockEvent } from "matchstick-as"
import { ethereum, Address, BigInt, Bytes } from "@graphprotocol/graph-ts"
import {
  AgreementEvent,
  RoleAdminChanged,
  RoleGranted,
  RoleRevoked
} from "../generated/FundsLock/FundsLock"

export function createAgreementEventEvent(
  seller: Address,
  buyer: Address,
  agreementId: BigInt,
  status: i32,
  timestamp: BigInt
): AgreementEvent {
  let mockEvent = newMockEvent()
  let agreementEventEvent = changetype<AgreementEvent>(newMockEvent())

  agreementEventEvent.parameters = new Array()

  agreementEventEvent.parameters.push(
    new ethereum.EventParam("seller", ethereum.Value.fromAddress(seller))
  )
  agreementEventEvent.parameters.push(
    new ethereum.EventParam("buyer", ethereum.Value.fromAddress(buyer))
  )
  agreementEventEvent.parameters.push(
    new ethereum.EventParam("id", ethereum.Value.fromUnsignedBigInt(agreementId))
  )
  agreementEventEvent.parameters.push(
    new ethereum.EventParam(
      "status",
      ethereum.Value.fromUnsignedBigInt(BigInt.fromI32(status))
    )
  )
  agreementEventEvent.parameters.push(
    new ethereum.EventParam(
      "timestamp",
      ethereum.Value.fromUnsignedBigInt(timestamp)
    )
  )

  return agreementEventEvent
}


