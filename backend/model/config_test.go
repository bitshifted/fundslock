// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"encoding/base64"
	"os"
	"testing"
)

const emptyGraphConfig = `[]`

func TestNewConfigurationLoader(t *testing.T) {
	loader := NewConfigurationLoader()
	if loader == nil {
		t.Fatal("NewConfigurationLoader returned nil")
	}
	if _, ok := loader.(*EnvironmentVariableConfigurationLoader); !ok {
		t.Errorf("Expected *EnvironmentVariableConfigurationLoader, got %T", loader)
	}
}

func TestLoad_MissingGraphConfigBase64(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	os.Setenv("GRAPH_URL", "http://localhost:8000")
	os.Setenv("GRAPH_API_KEY", "test-api-key")
	os.Setenv("GRAPH_CONFIG_BASE64", "")
	err := loader.Load()
	if err == nil {
		t.Fatal("Expected error for missing GRAPH_CONFIG_BASE64")
	}
	if want := "failed to decode GRAPH_CONFIG_BASE64: configuration variable GRAPH_CONFIG_BASE64 not set"; err.Error() != want {
		t.Errorf("Expected error %q, got %q", want, err.Error())
	}
}

func TestLoad_MissingJwtSecretKey(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	gc1 := `[{"chain":"mainnet","graphUrl":"http://subgraphs.api.url","apiKey":"graph-key"}]`
	os.Setenv("GRAPH_CONFIG_BASE64", base64.StdEncoding.EncodeToString([]byte(gc1)))
	os.Setenv("JWT_SECRET_KEY", "")
	err := loader.Load()
	if err == nil {
		t.Fatal("Expected error for missing JWT_SECRET_KEY")
	}
	if want := "configuration variable JWT_SECRET_KEY not set"; err.Error() != want {
		t.Errorf("Expected error %q, got %q", want, err.Error())
	}
}

func TestLoad_MissingCookieDomain(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	gc1 := `[{"chain":"mainnet","graphUrl":"http://subgraphs.api.url","apiKey":"graph-key"}]`
	os.Setenv("GRAPH_CONFIG_BASE64", base64.StdEncoding.EncodeToString([]byte(gc1)))
	os.Setenv("JWT_SECRET_KEY", "super-secret-jwt-key")
	os.Setenv("COOKIE_DOMAIN", "")
	err := loader.Load()
	if err == nil {
		t.Fatal("Expected error for missing COOKIE_DOMAIN")
	}
	expectedMsg := "configuration variable COOKIE_DOMAIN not set." +
		"Set it to the domain of your application (e.g., localhost for local development)"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestLoad_Success(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	os.Setenv("GRAPH_URL", "http://localhost:8000")
	os.Setenv("GRAPH_API_KEY", "test-graph-api-key")
	gc2 := `[{"chain":"mainnet","graphUrl":"http://subgraphs.api.url","apiKey":"graph-key"},` +
		`{"chain":"polygon","graphUrl":"http://polygon-subgraphs.api.url","apiKey":"polygon-key"}]`
	os.Setenv("GRAPH_CONFIG_BASE64", base64.StdEncoding.EncodeToString([]byte(gc2)))
	os.Setenv("JWT_SECRET_KEY", "super-secret-jwt-key")
	os.Setenv("ACCESS_TOKEN_DURATION", "1800")
	os.Setenv("REFRESH_TOKEN_DURATION", "604800")
	os.Setenv("SECURE_COOKIE", "true")
	os.Setenv("COOKIE_DOMAIN", "localhost")

	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(AppConfig.GraphConfig) != 2 {
		t.Fatalf("Expected 2 graph configs, got %d", len(AppConfig.GraphConfig))
	}
	if AppConfig.GraphConfig[0].Chain != "mainnet" {
		t.Errorf("Expected first config chain %q, got %q", "mainnet", AppConfig.GraphConfig[0].Chain)
	}
	if AppConfig.GraphConfig[1].Chain != "polygon" {
		t.Errorf("Expected second config chain %q, got %q", "polygon", AppConfig.GraphConfig[1].Chain)
	}
	if string(AppConfig.JwtSecretKey) != "super-secret-jwt-key" {
		t.Errorf("Expected JwtSecretKey %q, got %q", "super-secret-jwt-key", string(AppConfig.JwtSecretKey))
	}
	if AppConfig.AccessTokenDuration != 1800 {
		t.Errorf("Expected AccessTokenDuration %d, got %d", 1800, AppConfig.AccessTokenDuration)
	}
	if AppConfig.RefreshTokenDuration != 604800 {
		t.Errorf("Expected RefreshTokenDuration %d, got %d", 604800, AppConfig.RefreshTokenDuration)
	}
	if AppConfig.SecureCookie != true {
		t.Errorf("Expected SecureCookie true, got %t", AppConfig.SecureCookie)
	}
	if AppConfig.CookieDomain != "localhost" {
		t.Errorf("Expected CookieDomain %q, got %q", "localhost", AppConfig.CookieDomain)
	}
}

func setupBasicEnv() {
	os.Setenv("GRAPH_URL", "http://localhost:8000")
	os.Setenv("GRAPH_API_KEY", "test-api-key")
	gc := emptyGraphConfig
	os.Setenv("GRAPH_CONFIG_BASE64", base64.StdEncoding.EncodeToString([]byte(gc)))
	os.Setenv("JWT_SECRET_KEY", "secret")
	os.Setenv("COOKIE_DOMAIN", "localhost")
}

func TestLoad_DefaultAccessTokenDuration(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if AppConfig.AccessTokenDuration != defaultAccessTokenDuration {
		t.Errorf("Expected default AccessTokenDuration %d, got %d", defaultAccessTokenDuration, AppConfig.AccessTokenDuration)
	}
}

func TestLoad_ZeroAccessTokenDuration(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	os.Setenv("ACCESS_TOKEN_DURATION", "0")
	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if AppConfig.AccessTokenDuration != defaultAccessTokenDuration {
		t.Errorf("Expected default AccessTokenDuration %d when 0 is provided, got %d",
			defaultAccessTokenDuration, AppConfig.AccessTokenDuration)
	}
}

func TestLoad_DefaultRefreshTokenDuration(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if AppConfig.RefreshTokenDuration != defaultRefreshTokenDuration {
		t.Errorf("Expected default RefreshTokenDuration %d, got %d", defaultRefreshTokenDuration, AppConfig.RefreshTokenDuration)
	}
}

func TestLoad_ZeroRefreshTokenDuration(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	os.Setenv("REFRESH_TOKEN_DURATION", "0")
	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if AppConfig.RefreshTokenDuration != defaultRefreshTokenDuration {
		t.Errorf("Expected default RefreshTokenDuration %d when 0 is provided, got %d",
			defaultRefreshTokenDuration, AppConfig.RefreshTokenDuration)
	}
}

func TestLoad_InvalidSecureCookieDefaultsToTrue(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	os.Setenv("SECURE_COOKIE", "not-a-boolean")
	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if AppConfig.SecureCookie != true {
		t.Errorf("Expected SecureCookie to default to true for invalid value, got %t", AppConfig.SecureCookie)
	}
}

func TestLoad_EmptySecureCookieDefaultsToTrue(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if AppConfig.SecureCookie != true {
		t.Errorf("Expected SecureCookie to default to true when not set, got %t", AppConfig.SecureCookie)
	}
}

func TestLoad_SecureCookieFalse(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	os.Setenv("SECURE_COOKIE", "false")
	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if AppConfig.SecureCookie != false {
		t.Errorf("Expected SecureCookie false, got %t", AppConfig.SecureCookie)
	}
}

func TestLoad_SecureCookieTrue(t *testing.T) {
	loader := NewConfigurationLoader()
	os.Clearenv()
	setupBasicEnv()
	os.Setenv("SECURE_COOKIE", "true")
	err := loader.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if AppConfig.SecureCookie != true {
		t.Errorf("Expected SecureCookie true, got %t", AppConfig.SecureCookie)
	}
}

func TestDecodeGraphConfig_EmptyEnv(t *testing.T) {
	os.Clearenv()
	_, err := decodeGraphConfig()
	if err == nil {
		t.Fatal("Expected error for empty GRAPH_CONFIG_BASE64")
	}
	if want := "configuration variable GRAPH_CONFIG_BASE64 not set"; err.Error() != want {
		t.Errorf("Expected error %q, got %q", want, err.Error())
	}
}

func TestDecodeGraphConfig_InvalidBase64(t *testing.T) {
	os.Clearenv()
	os.Setenv("GRAPH_CONFIG_BASE64", "not-valid-base64!!!")
	_, err := decodeGraphConfig()
	if err == nil {
		t.Fatal("Expected error for invalid base64")
	}
}

func TestDecodeGraphConfig_InvalidJson(t *testing.T) {
	os.Clearenv()
	os.Setenv("GRAPH_CONFIG_BASE64", base64.StdEncoding.EncodeToString([]byte(`{not valid json}`)))
	_, err := decodeGraphConfig()
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

func TestDecodeGraphConfig_Valid(t *testing.T) {
	os.Clearenv()
	config := `{"chain":"ethereum","graphUrl":"http://api.example.com","apiKey":"key123"}`
	gc := "[" + config + "]"
	os.Setenv("GRAPH_CONFIG_BASE64", base64.StdEncoding.EncodeToString([]byte(gc)))
	result, err := decodeGraphConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 config, got %d", len(result))
	}
	if result[0].Chain != "ethereum" {
		t.Errorf("Expected chain %q, got %q", "ethereum", result[0].Chain)
	}
	if result[0].GraphUrl != "http://api.example.com" {
		t.Errorf("Expected graphUrl %q, got %q", "http://api.example.com", result[0].GraphUrl)
	}
	if result[0].ApiKey != "key123" {
		t.Errorf("Expected apiKey %q, got %q", "key123", result[0].ApiKey)
	}
}

func TestDecodeGraphConfig_EmptyArray(t *testing.T) {
	os.Clearenv()
	os.Setenv("GRAPH_CONFIG_BASE64", base64.StdEncoding.EncodeToString([]byte(emptyGraphConfig)))
	result, err := decodeGraphConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 configs, got %d", len(result))
	}
}

func TestDecodeGraphConfig_MultipleConfigs(t *testing.T) {
	os.Clearenv()
	c1 := `{"chain":"mainnet","graphUrl":"http://main.api","apiKey":"key1"}`
	c2 := `{"chain":"goerli","graphUrl":"http://goerli.api","apiKey":"key2"}`
	c3 := `{"chain":"sepolia","graphUrl":"http://sepolia.api","apiKey":"key3"}`
	configs := "[" + c1 + "," + c2 + "," + c3 + "]"
	os.Setenv("GRAPH_CONFIG_BASE64", base64.StdEncoding.EncodeToString([]byte(configs)))
	result, err := decodeGraphConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("Expected 3 configs, got %d", len(result))
	}
	expectedChains := []string{"mainnet", "goerli", "sepolia"}
	for i, expected := range expectedChains {
		if result[i].Chain != expected {
			t.Errorf("Expected config[%d].Chain %q, got %q", i, expected, result[i].Chain)
		}
	}
}
