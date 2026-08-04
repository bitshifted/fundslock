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
	walletAddress, err := auth.VerifyMessage(verificationRequest)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to verify SIWE message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokens, err := auth.GenerateTokens(walletAddress)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to generate tokens:")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.Info().Msgf("Generated tokens for wallet %s: ", walletAddress)
	// set refresh token as httpOnly cookie
	//nolint:gosec // G124 ignoring gosec warning for cookie config
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   model.AppConfig.SecureCookie,
		Path:     "/api/v1/auth/refresh",
		SameSite: http.SameSiteLaxMode,
	})
	accessTokenResponse := model.AccessTokenResponse{
		AccessToken: tokens.AccessToken,
	}
	//nolint:gosec // G117 ignoring gosec warning for using access_token field
	err = json.NewEncoder(w).Encode(accessTokenResponse)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to encode access token response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func refreshAccessToken(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info().Msg("Refreshing access token...")
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to retrieve refresh token cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	refreshToken := cookie.Value
	log.Logger.Info().Msgf("Validating refresh token: %s", refreshToken)
	claims, err := auth.ValidateToken(refreshToken)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to validate refresh token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokens, err := auth.GenerateTokens(claims.WalletAddress)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to generate tokens:")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.Info().Msgf("Generated new tokens for wallet %s: ", claims.WalletAddress)
	// set new refresh token as httpOnly cookie
	//nolint:gosec // G124 ignoring gosec warning for cookie config
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   model.AppConfig.SecureCookie,
		Path:     "/api/v1/auth/refresh",
		SameSite: http.SameSiteLaxMode,
	})
	accessTokenResponse := model.AccessTokenResponse{
		AccessToken: tokens.AccessToken,
	}
	//nolint:gosec // G117 ignoring gosec warning for using access_token field
	err = json.NewEncoder(w).Encode(accessTokenResponse)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to encode access token response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
