// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package model

import "time"

type AgreementStatus int8

// mapping agreement status to smart contract values
const (
	Created                AgreementStatus = 0
	Funded                 AgreementStatus = 1
	SellerAccepted         AgreementStatus = 2
	SellerRequestedRelease AgreementStatus = 3
	BuyerApprovedRelease   AgreementStatus = 4
	Released               AgreementStatus = 5
	Canceled               AgreementStatus = 6
)

var AgreementStatusMap = map[AgreementStatus]string{
	Created:                "created",
	Funded:                 "funded",
	SellerAccepted:         "seller_accepted",
	SellerRequestedRelease: "seller_requested_release",
	BuyerApprovedRelease:   "buyer_approved_release",
	Released:               "released",
	Canceled:               "canceled",
}

func (s AgreementStatus) String() string {
	return AgreementStatusMap[s]
}

type AgreementData struct {
	AgreementId string    `json:"agreement_id"`
	Seller      string    `json:"seller"`
	Buyer       string    `json:"buyer"`
	Amount      string    `json:"amount"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}
