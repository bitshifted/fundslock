// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"crypto/ecdsa"
	"fmt"
	"strings"
	"testing"
	"time"

	"bitshifted/fundslock-be/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spruceid/siwe-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type walletTestSuite struct {
	suite.Suite
	testKey     *ecdsa.PrivateKey
	testAddress string
}

func (s *walletTestSuite) SetupSuite() {
	var err error
	s.testKey, err = crypto.GenerateKey()
	assert.NoError(s.T(), err)
	s.testAddress = crypto.PubkeyToAddress(s.testKey.PublicKey).Hex()
}

func (s *walletTestSuite) SetupTest() {
	nonceStore = make(map[string]time.Time)
}

func (s *walletTestSuite) buildSIWEMessage(nonce string) string {
	return fmt.Sprintf(`%s wants you to sign in with your Ethereum account:
%s


URI: https://example.com
Version: 1
Chain ID: 1
Nonce: %s
Issued At: 2026-08-04T13:45:28Z`, "example.com", s.testAddress, nonce)
}

func (s *walletTestSuite) signMessage(message string) (string, error) {
	hash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))
	sig, err := crypto.Sign(hash.Bytes(), s.testKey)
	if err != nil {
		return "", err
	}
	sig[64] += 27
	return "0x" + common.Bytes2Hex(sig), nil
}

func (s *walletTestSuite) TestVerifyMessageHappyPath() {
	nonce := "abcdefgh12345678"
	message := s.buildSIWEMessage(nonce)
	signature, err := s.signMessage(message)
	assert.NoError(s.T(), err)

	nonceStore[nonce] = time.Now().Add(1 * time.Minute)

	addr, err := VerifyMessage(model.SIWEVerificationRequest{
		Message:   message,
		Signature: signature,
		Nonce:     nonce,
	})
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), addr)
	_, exists := nonceStore[nonce]
	assert.False(s.T(), exists, "nonce should be removed after verification")
}

func (s *walletTestSuite) TestVerifyMessageInvalidMessage() {
	nonceStore["valid-nonce"] = time.Now().Add(1 * time.Minute)

	_, err := VerifyMessage(model.SIWEVerificationRequest{
		Message:   "this is not a valid SIWE message at all",
		Signature: "0xdeadbeef",
		Nonce:     "valid-nonce",
	})
	assert.Error(s.T(), err)
	_, exists := nonceStore["valid-nonce"]
	assert.True(s.T(), exists, "nonce should not be removed on parse failure")
}

func (s *walletTestSuite) TestVerifyMessageExpiredNonce() {
	nonce := "expirednonce12345"
	message := s.buildSIWEMessage(nonce)

	nonceStore[nonce] = time.Now().Add(-2 * time.Minute)

	_, err := VerifyMessage(model.SIWEVerificationRequest{
		Message:   message,
		Signature: "0xdeadbeef",
		Nonce:     nonce,
	})
	assert.Error(s.T(), err)
	assert.EqualError(s.T(), err, "invalid or expired nonce")
}

func (s *walletTestSuite) TestVerifyMessageMissingNonce() {
	message := s.buildSIWEMessage("nonexistentnonce")

	_, err := VerifyMessage(model.SIWEVerificationRequest{
		Message:   message,
		Signature: "0xdeadbeef",
		Nonce:     "nonexistentnonce",
	})
	assert.Error(s.T(), err)
	assert.EqualError(s.T(), err, "invalid or expired nonce")
}

func (s *walletTestSuite) TestVerifyMessageInvalidSignature() {
	nonce := "sigtestnonce12345"
	message := s.buildSIWEMessage(nonce)

	nonceStore[nonce] = time.Now().Add(1 * time.Minute)

	_, err := VerifyMessage(model.SIWEVerificationRequest{
		Message:   message,
		Signature: "0x" + strings.Repeat("0", 130),
		Nonce:     nonce,
	})
	assert.Error(s.T(), err)
}

func (s *walletTestSuite) TestVerifyMessageWrongWalletSignature() {
	nonce := "wrongwalletnonce1"
	message := s.buildSIWEMessage(nonce)

	wrongKey, err := crypto.GenerateKey()
	assert.NoError(s.T(), err)
	wrongSig, err := s.signMessageWith(wrongKey, message)
	assert.NoError(s.T(), err)

	nonceStore[nonce] = time.Now().Add(1 * time.Minute)

	_, err = VerifyMessage(model.SIWEVerificationRequest{
		Message:   message,
		Signature: wrongSig,
		Nonce:     nonce,
	})
	assert.Error(s.T(), err)
}

func (s *walletTestSuite) signMessageWith(key *ecdsa.PrivateKey, message string) (string, error) {
	hash := crypto.Keccak256Hash(fmt.Appendf(nil, "\x19Ethereum Signed Message:\n%d%s", len(message), message))
	sig, err := crypto.Sign(hash.Bytes(), key)
	if err != nil {
		return "", err
	}
	sig[64] += 27
	return "0x" + common.Bytes2Hex(sig), nil
}

func (s *walletTestSuite) TestVerifyMessageNonAlphanumericNonce() {
	nonceStore["has-hyphen-nonce"] = time.Now().Add(1 * time.Minute)

	message := s.buildSIWEMessage("has-hyphen-nonce")

	_, err := VerifyMessage(model.SIWEVerificationRequest{
		Message:   message,
		Signature: "0xdeadbeef",
		Nonce:     "has-hyphen-nonce",
	})
	assert.Error(s.T(), err)
}

func (s *walletTestSuite) TestVerifyMessageNonceTooShort() {
	nonceStore["shortnonce"] = time.Now().Add(1 * time.Minute)

	message := s.buildSIWEMessage("shortnonce")

	_, err := VerifyMessage(model.SIWEVerificationRequest{
		Message:   message,
		Signature: "0x" + strings.Repeat("1", 130),
		Nonce:     "shortnonce",
	})
	assert.Error(s.T(), err)
}

func (s *walletTestSuite) TestVerifyMessageTrailingWhitespaceInInitMessage() {
	nonce := "trimtestnonce1234"
	initMsg, err := siwe.InitMessage(
		"localhost:3000",
		s.testAddress,
		"http://localhost:3000",
		nonce,
		nil,
	)
	assert.NoError(s.T(), err)

	cleaned := strings.TrimSpace(initMsg.String())
	signature, err := s.signMessage(cleaned)
	assert.NoError(s.T(), err)

	nonceStore[nonce] = time.Now().Add(1 * time.Minute)

	addr, err := VerifyMessage(model.SIWEVerificationRequest{
		Message:   cleaned,
		Signature: signature,
		Nonce:     nonce,
	})
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), addr)
}

func TestWalletTestSuite(t *testing.T) {
	suite.Run(t, new(walletTestSuite))
}
