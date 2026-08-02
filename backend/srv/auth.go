// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package srv

import (
	"bitshifted/fundslock-be/auth"
	"bitshifted/fundslock-be/log"
	"bitshifted/fundslock-be/model"
	"encoding/json"
	"net/http"
)

func createNonce(w http.ResponseWriter, r *http.Request) {
	nonce := auth.GenerateNonce()
	err := json.NewEncoder(w).Encode(nonce)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func verifySIWEMessage(w http.ResponseWriter, r *http.Request) {
	verificationRequest := model.SIWEVerificationRequest{}
	err := json.NewDecoder(r.Body).Decode(&verificationRequest)
	log.Logger.Info().Msgf("Received SIWE verification request: %+v", verificationRequest)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to decode SIWE verification request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokens, err := auth.VerifyMessage(verificationRequest)
	log.Logger.Info().Msgf("Generated tokens: %+v", tokens)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to verify SIWE message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to encode tokens")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
