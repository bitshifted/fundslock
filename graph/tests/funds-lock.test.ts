import {
  assert,
  describe,
  test,
  clearStore,
  beforeAll,
  afterAll
} from "matchstick-as/assembly/index"
import { Address, BigInt, Bytes } from "@graphprotocol/graph-ts"
import { AgreementLog } from "../generated/schema"
import { handleAgreementEvent } from "../src/funds-lock"
import { createAgreementEventEvent } from "./funds-lock-utils"

// Tests structure (matchstick-as >=0.5.0)
// https://thegraph.com/docs/en/subgraphs/developing/creating/unit-testing-framework/#tests-structure

describe("Describe entity assertions", () => {
  // beforeAll(() => {
  //   let seller = Address.fromString(
  //     "0x0000000000000000000000000000000000000001"
  //   )
  //   let buyer = Address.fromString("0x0000000000000000000000000000000000000002")
  //   let amount = BigInt.fromI32(234)
  //   let status = 123
  //   let timestamp = BigInt.fromI32(234)
  //   let newAgreementEventEvent = createAgreementEventEvent(
  //     seller,
  //     buyer,
  //     amount,
  //     status,
  //     timestamp
  //   )
  //   handleAgreementEvent(newAgreementEventEvent)
  // })

  // afterAll(() => {
  //   clearStore()
  // })

  // For more test scenarios, see:

  test("AgreementEvent created and stored", () => {
    // assert.entityCount("AgreementLog", 1)

     let seller = Address.fromString(
      "0x0000000000000000000000000000000000000001"
    )
    let buyer = Address.fromString("0x0000000000000000000000000000000000000002")
    let agreementID = BigInt.fromI32(234)
    let status = 123
    let timestamp = BigInt.fromI32(234)

    let event = createAgreementEventEvent(
       seller,
      buyer,
      agreementID,
      status,
      timestamp
    )
    handleAgreementEvent(event)

    // assert.fieldEquals(
    //   "AgreementLog",
    //   "123",
    //   "agreement_id",
    //   agreementID.toString()
    // )
    

    // 0xa16081f360e3847006db660bae1c6d1b2e17ec2a is the default address used in newMockEvent() function
    // Log index is 0 by default in mock events
    // let id = "wrong-id"
    // assert.fieldEquals(
    //   "AgreementLog",
    //   id,
    //   "seller",
    //   "0x0000000000000000000000000000000000000001"
    // )
    // assert.fieldEquals(
    //   "AgreementLog",
    //   id,
    //   "buyer",
    //   "0x0000000000000000000000000000000000000002"
    // )
    // assert.fieldEquals(
    //   "AgreementLog",
    //   id,
    //   "status",
    //   "123"
    // )
    // assert.fieldEquals(
    //   "AgreementLog",
    //   id,
    //   "timestamp",
    //   "234"
    // )
  })
})
