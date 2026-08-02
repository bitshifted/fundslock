package auth

import (
	"bitshifted/fundslock-be/log"
	"bitshifted/fundslock-be/model"
	"fmt"
	"time"

	"github.com/spruceid/siwe-go"
)

func VerifyMessage(requset model.SIWEVerificationRequest) (*model.TokenPair, error) {
	siweMsg, err := siwe.ParseMessage(requset.Message)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to parse SIWE message")
		return nil, err
	}
	fmt.Println("Parsed SIWE message:", siweMsg)
	nonceExpirationTime := nonceStore[requset.Nonce]
	if nonceExpirationTime.Before(time.Now()) {
		log.Logger.Error().Msg("Invalid or expired nonce")
		return nil, fmt.Errorf("invalid or expired nonce")
	}
	removeNonce(requset.Nonce)

	_, err = siweMsg.Verify(requset.Signature, nil, &requset.Nonce, nil)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to verify SIWE message")
		return nil, err
	}
	walletAddress := siweMsg.GetAddress().Hex()
	tokens, err := GenerateTokens(walletAddress)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to generate tokens:")
		return nil, err
	}
	log.Logger.Info().Msgf("Generated tokens for wallet %s: ", walletAddress)
	return tokens, nil

}
