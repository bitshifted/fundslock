package auth

import (
	"bitshifted/fundslock-be/log"
	"bitshifted/fundslock-be/model"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	nonceStore = make(map[string]time.Time)
)

type Claims struct {
	WalletAddress string `json:"wallet_address"`
	jwt.RegisteredClaims
}

func nonceValue(nonceId string) time.Time {
	expiration, exists := nonceStore[nonceId]
	if !exists {
		return time.Unix(0, 0)
	}
	return expiration
}

func removeNonce(nonce string) {
	delete(nonceStore, nonce)
}

func GenerateNonce() model.Nonce {
	log.Logger.Info().Msg("Creating nonce")
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	nonce := model.Nonce{
		Nonce: hex.EncodeToString(b),
	}

	nonceStore[nonce.Nonce] = time.Now().Add(1 * time.Minute) // Set expiration for 1 minute
	log.Logger.Info().Str("nonce", nonce.Nonce).Msgf("Nonce created and stored with expiration %s",
		nonceStore[nonce.Nonce].Format(time.RFC3339))
	return nonce
}

func GenerateTokens(walletAddress string) (*model.TokenPair, error) {
	if len(model.AppConfig.JwtSecretKey) == 0 {
		return nil, errors.New("JWT secret key is not set in environment variables")
	}
	now := time.Now()
	accessTokenClaims := &Claims{
		WalletAddress: walletAddress,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)), // Access token valid for 15 minutes
			Subject:   walletAddress,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccess, err := accessToken.SignedString(model.AppConfig.JwtSecretKey)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to sign access token")
		return nil, err
	}

	//  Refresh Token Claims (7 days)
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(7 * time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(now),
		Subject:   walletAddress,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString(model.AppConfig.JwtSecretKey)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to sign refresh token")
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  signedAccess,
		RefreshToken: signedRefresh,
	}, nil
}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return model.AppConfig.JwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
