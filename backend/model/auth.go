// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package model

type Nonce struct {
	Nonce string `json:"nonce"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type SIWEVerificationRequest struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
}
