// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"bitshifted/fundslock-be/log"
	"errors"
	"os"
	"strconv"
)

type ConfigurationVariables struct {
	GraphUrl             string
	GraphApiKey          string
	JwtSecretKey         []byte
	AccessTokenDuration  int64 // in seconds
	RefreshTokenDuration int64 // in seconds
	SecureCookie         bool  // for local development, set to false to allow non-secure cookies
	CookieDomain         string
}

const (
	defaultAccessTokenDuration  = 15 * 60          // 15 minutes
	defaultRefreshTokenDuration = 7 * 24 * 60 * 60 // 7 days
)

type ConfifgurationLoader interface {
	Load() error
}

type EnvironmentVariableConfigurationLoader struct {
	ConfifgurationLoader
}

func NewConfigurationLoader() ConfifgurationLoader {
	return &EnvironmentVariableConfigurationLoader{}
}

var AppConfig = ConfigurationVariables{}

func (e *EnvironmentVariableConfigurationLoader) Load() error {
	graphUrl := os.Getenv("GRAPH_URL")
	if graphUrl == "" {
		return errors.New("configuration variable GRAPH_URL not set")
	}
	graphApiKey := os.Getenv("GRAPH_API_KEY")
	if graphApiKey == "" {
		return errors.New("configuration variable GRAPH_API_KEY not set")
	}
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return errors.New("configuration variable JWT_SECRET_KEY not set")
	}
	accessTokenDuration, _ := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_DURATION"), 10, 64)
	if accessTokenDuration == 0 {
		accessTokenDuration = defaultAccessTokenDuration
	}
	refreshTokenDuration, _ := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_DURATION"), 10, 64)
	if refreshTokenDuration == 0 {
		refreshTokenDuration = defaultRefreshTokenDuration
	}
	secureCookie, err := strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
	if err != nil {
		log.Logger.Warn().Msg("SECURE_COOKIE is not set or can't be read. Using default value: true")
		secureCookie = true
	}
	cookieDDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDDomain == "" {
		return errors.New("configuration variable COOKIE_DOMAIN not set. Set it to the domain of your application (e.g., localhost for local development)")
	}
	AppConfig.GraphUrl = graphUrl
	AppConfig.GraphApiKey = graphApiKey
	AppConfig.JwtSecretKey = []byte(jwtSecretKey)
	AppConfig.AccessTokenDuration = accessTokenDuration
	AppConfig.RefreshTokenDuration = refreshTokenDuration
	AppConfig.SecureCookie = secureCookie
	AppConfig.CookieDomain = cookieDDomain
	return nil
}
