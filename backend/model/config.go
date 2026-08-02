// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"errors"
	"os"
)

type ConfigurationVariables struct {
	GraphUrl     string
	GraphApiKey  string
	JwtSecretKey []byte
}

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
	AppConfig.GraphUrl = graphUrl
	AppConfig.GraphApiKey = graphApiKey
	AppConfig.JwtSecretKey = []byte(jwtSecretKey)
	return nil
}
