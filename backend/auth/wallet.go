// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"bitshifted/fundslock-be/log"
	"bitshifted/fundslock-be/model"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spruceid/siwe-go"
)

// VerifyMessage verifies the SIWE message and returns wallet address if the verification is successful.
func VerifyMessage(requset model.SIWEVerificationRequest) (*model.SessionData, error) {
	siweMsg, err := siwe.ParseMessage(requset.Message)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to parse SIWE message")
		return nil, err
	}
	log.Logger.Debug().Msgf("Parsed SIWE message: %s", siweMsg)
	nonce := parseNonce(strings.ToLower(requset.Message))
	log.Logger.Debug().Msgf("Parsed nonce: %s", nonce)
	nonceExpirationTime, ok := nonceStore[nonce]
	if !ok {
		log.Logger.Error().Msgf("Nonce %s not found", nonce)
		return nil, fmt.Errorf("nonce not found")
	} else {
		log.Logger.Info().Msgf("Nonce found with expiration time: %s", nonceExpirationTime.Format(time.RFC3339))
	}
	log.Logger.Info().Msgf("Nonce expiration time: %s", nonceExpirationTime.Format(time.RFC3339))
	if nonceExpirationTime.Before(time.Now()) {
		log.Logger.Error().Msg("Invalid or expired nonce")
		return nil, fmt.Errorf("invalid or expired nonce")
	}
	removeNonce(nonce)

	_, err = siweMsg.Verify(requset.Signature, nil, &nonce, nil)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to verify SIWE message")
		return nil, err
	}
	return &model.SessionData{
		WalletAddress: siweMsg.GetAddress().Hex(),
		ChainId:       uint32(siweMsg.GetChainID()),
	}, nil
}

func parseNonce(input string) string {
	log.Logger.Info().Msgf("Parsing message: %s", input)
	re := regexp.MustCompile(`\snonce:\s*([0-9a-z]+)\s`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		log.Logger.Error().Msg("Failed to parse nonce")
		return ""
	}
	log.Logger.Info().Msgf("Parsed nonce: %s", matches[0])
	return strings.TrimSpace(matches[1])
}
