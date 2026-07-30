package model

import (
	"errors"
	"os"
)

type ConfigurationVariables struct {
	GraphUrl    string
	GraphApiKey string
}

type ConfifgurationLoader interface {
	Load() (*ConfigurationVariables, error)
}

type EnvironmentVariableConfigurationLoader struct {
	ConfifgurationLoader
}

func NewConfigurationLoader() ConfifgurationLoader {
	return &EnvironmentVariableConfigurationLoader{}
}

func (e *EnvironmentVariableConfigurationLoader) Load() (*ConfigurationVariables, error) {
	graphUrl := os.Getenv("GRAPH_URL")
	if graphUrl == "" {
		return nil, errors.New("Configuration variable GRAPH_URL not set")
	}
	graphApiKey := os.Getenv("GRAPH_API_KEY")
	if graphApiKey == "" {
		return nil, errors.New("Configuration variable GRAPH_API_KEY not set")
	}

	return &ConfigurationVariables{
		GraphUrl:    graphUrl,
		GraphApiKey: graphApiKey,
	}, nil

}
