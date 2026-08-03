// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"bitshifted/fundslock-be/log"
	"bitshifted/fundslock-be/model"
	"fmt"
	"time"

	"github.com/spruceid/siwe-go"
)

// VerifyMessage verifies the SIWE message and returns wallet address if the verification is successful.
func VerifyMessage(requset model.SIWEVerificationRequest) (string, error) {
	siweMsg, err := siwe.ParseMessage(requset.Message)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to parse SIWE message")
		return "", err
	}
	fmt.Println("Parsed SIWE message:", siweMsg)
	nonceExpirationTime := nonceStore[requset.Nonce]
	if nonceExpirationTime.Before(time.Now()) {
		log.Logger.Error().Msg("Invalid or expired nonce")
		return "", fmt.Errorf("invalid or expired nonce")
	}
	removeNonce(requset.Nonce)

	_, err = siweMsg.Verify(requset.Signature, nil, &requset.Nonce, nil)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to verify SIWE message")
		return "", err
	}
	return siweMsg.GetAddress().Hex(), nil
}
