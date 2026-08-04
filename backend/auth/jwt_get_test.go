// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"bitshifted/fundslock-be/model"

	"github.com/golang-jwt/jwt/v5"
)

// Test suite for JWT functions

type jwtTestSuite struct {
	suite.Suite
}

func (s *jwtTestSuite) SetupTest() {
	// Reset nonce store and config before each test
	nonceStore = make(map[string]time.Time)
	model.AppConfig = model.ConfigurationVariables{
		JwtSecretKey:         []byte("test-secret-key"),
		AccessTokenDuration:  60,  // 1 minute
		RefreshTokenDuration: 120, // 2 minutes
	}
}

func (s *jwtTestSuite) TestGenerateNonce() {
	nonce := GenerateNonce()
	assert.NotEmpty(s.T(), nonce.Nonce)
	// nonce should be stored with expiration
	exp, ok := nonceStore[nonce.Nonce]
	assert.True(s.T(), ok)
	assert.True(s.T(), exp.After(time.Now()))
}

func (s *jwtTestSuite) TestGenerateTokensSuccess() {
	pair, err := GenerateTokens("0xabc123")
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), pair.AccessToken)
	assert.NotEmpty(s.T(), pair.RefreshToken)

	// Validate access token
	claims, err := ValidateToken(pair.AccessToken)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "0xabc123", claims.WalletAddress)
}

func (s *jwtTestSuite) TestGenerateTokensMissingSecret() {
	// Clear secret key
	model.AppConfig.JwtSecretKey = nil
	_, err := GenerateTokens("0xabc123")
	assert.Error(s.T(), err)
	assert.EqualError(s.T(), err, "JWT secret key is not set in environment variables")
}

func (s *jwtTestSuite) TestValidateTokenInvalid() {
	// Token signed with different key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{WalletAddress: "0xabc"})
	signed, _ := token.SignedString([]byte("other-key"))
	_, err := ValidateToken(signed)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "signature is invalid")
}

func TestJwtTestSuite(t *testing.T) {
	suite.Run(t, new(jwtTestSuite))
}
