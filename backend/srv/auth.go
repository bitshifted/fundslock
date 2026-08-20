// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package srv

import (
	"bitshifted/fundslock-be/auth"
	"bitshifted/fundslock-be/log"
	"bitshifted/fundslock-be/model"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

const (
	cookieExpirationPeriod = 6
)

var tokenAuth *jwtauth.JWTAuth

func jwtInit() {
	tokenAuth = jwtauth.New("HS256", (model.AppConfig.JwtSecretKey), nil)
}

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
	sessionData, err := auth.VerifyMessage(verificationRequest)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to verify SIWE message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokens, err := auth.GenerateTokens(sessionData.WalletAddress, sessionData.ChainId)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to generate tokens:")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.Info().Msgf("Generated tokens for wallet %s: ", sessionData.WalletAddress)
	// set refresh token as httpOnly cookie
	//nolint:gosec // G124 ignoring gosec warning for cookie config
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   model.AppConfig.SecureCookie,
		Path:     "/api/v1/auth/refresh",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(model.AppConfig.RefreshTokenDuration) * time.Second),
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
	if claims.ExpiresAt == nil {
		log.Logger.Error().Msg("Refresh token has no expiration time")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokens, err := auth.GenerateTokens(claims.WalletAddress, claims.ChainId)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to generate tokens:")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.Info().Msgf("Generated new tokens for wallet %s: ", claims.WalletAddress)
	log.Logger.Debug().Msgf("Refresh token expiration time: %s", claims.ExpiresAt.Format(time.RFC3339))
	// set new refresh token as httpOnly cookie, only if it expires in 6 hours
	if claims.ExpiresAt.Before(time.Now().Add(cookieExpirationPeriod * time.Hour)) {
		log.Logger.Debug().Msg("Cookie expires less than 6 hours from now, generating new cookie")
		//nolint:gosec // G124 ignoring gosec warning for cookie config
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			HttpOnly: true,
			Secure:   model.AppConfig.SecureCookie,
			Path:     "/api/v1/auth/refresh",
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(time.Duration(model.AppConfig.RefreshTokenDuration) * time.Second),
		})
	} else {
		log.Logger.Debug().Msg("Cookie expires more than 6 hours from now, using existing cookie")
	}

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

func getSession(w http.ResponseWriter, r *http.Request) {
	log.Logger.Debug().Msg("Retrieveing session data...")
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to retrieve session data")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	walletAddress, ok := claims["wallet_address"].(string)
	if !ok {
		log.Logger.Error().Msg("Failed to retrieve wallet address from JWT token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	chainId, ok := claims["chain_id"].(float64)
	if !ok {
		log.Logger.Error().Msg("Failed to retrieve chain ID from JWT token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sessionData := model.SessionData{
		WalletAddress: walletAddress,
		ChainId:       int(chainId),
	}
	err = json.NewEncoder(w).Encode(sessionData)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to encode session data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
