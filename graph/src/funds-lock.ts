import { BigInt } from "@graphprotocol/graph-ts"
import {
  AgreementEvent as AgreementEventEvent,
  AgreementCreated as AgreementCreatedEvent,
} from "../generated/FundsLock/FundsLock"
import {
  AgreementLog
} from "../generated/schema"

export function handleAgreementEvent(event: AgreementEventEvent): void {
  let entity = new AgreementLog(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.seller = event.params.seller
  entity.buyer = event.params.buyer
  entity.agreement_id = event.params.id
  entity.status = event.params.status
  entity.amount = new BigInt(0)
  entity.timestamp = event.params.timestamp

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleAgreementCreated(event: AgreementCreatedEvent) : void {
  let entity = new AgreementLog(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.seller = event.params.seller
  entity.buyer = event.params.buyer
  entity.agreement_id = event.params.id
  entity.amount = event.params.amount
  entity.status = 0
  entity.timestamp = event.params.timestamp


  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}
