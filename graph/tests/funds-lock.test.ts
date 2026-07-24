import {
  assert,
  describe,
  test,
  clearStore,
  beforeAll,
  afterAll
} from "matchstick-as/assembly/index"
import { Address, BigInt, Bytes } from "@graphprotocol/graph-ts"
import { AgreementEvent } from "../generated/schema"
import { AgreementEvent as AgreementEventEvent } from "../generated/FundsLock/FundsLock"
import { handleAgreementEvent } from "../src/funds-lock"
import { createAgreementEventEvent } from "./funds-lock-utils"

// Tests structure (matchstick-as >=0.5.0)
// https://thegraph.com/docs/en/subgraphs/developing/creating/unit-testing-framework/#tests-structure

describe("Describe entity assertions", () => {
  beforeAll(() => {
    let seller = Address.fromString(
      "0x0000000000000000000000000000000000000001"
    )
    let buyer = Address.fromString("0x0000000000000000000000000000000000000001")
    let amount = BigInt.fromI32(234)
    let status = 123
    let timestamp = BigInt.fromI32(234)
    let newAgreementEventEvent = createAgreementEventEvent(
      seller,
      buyer,
      amount,
      status,
      timestamp
    )
    handleAgreementEvent(newAgreementEventEvent)
  })

  afterAll(() => {
    clearStore()
  })

  // For more test scenarios, see:

  test("AgreementEvent created and stored", () => {
    assert.entityCount("AgreementEvent", 1)

    // 0xa16081f360e3847006db660bae1c6d1b2e17ec2a is the default address used in newMockEvent() function
    // assert.fieldEquals(
    //   "AgreementEvent",
    //   "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
    //   "seller",
    //   "0x0000000000000000000000000000000000000001"
    // )
    // assert.fieldEquals(
    //   "AgreementEvent",
    //   "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
    //   "buyer",
    //   "0x0000000000000000000000000000000000000001"
    // )
    // assert.fieldEquals(
    //   "AgreementEvent",
    //   "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
    //   "amount",
    //   "234"
    // )
    // assert.fieldEquals(
    //   "AgreementEvent",
    //   "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
    //   "status",
    //   "123"
    // )
    // assert.fieldEquals(
    //   "AgreementEvent",
    //   "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
    //   "timestamp",
    //   "234"
    // )

    // More assert options:
    // https://thegraph.com/docs/en/subgraphs/developing/creating/unit-testing-framework/#asserts
  })
})
