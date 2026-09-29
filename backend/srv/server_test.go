// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package srv

import (
	"bitshifted/fundslock-be/common"
	"encoding/base64"
	"encoding/json"
	"os"
	"sync"
	"testing"

	"bitshifted/fundslock-be/model"
)

func TestStartLambdaAdapterInvoked(t *testing.T) {
	if common.IsLambdaEnvironment() {
		t.Skip("already running in Lambda environment")
	}

	const lambdaEnvVar = "AWS_LAMBDA_FUNCTION_NAME"
	err := os.Setenv(lambdaEnvVar, "test-lambda-function")
	if err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	defer os.Unsetenv(lambdaEnvVar)

	graphConfig := []model.GraphConfig{
		{Chain: "test", GraphUrl: "http://localhost:8080", ApiKey: "test-key"},
	}
	marshaled, _ := json.Marshal(graphConfig) //nolint:gosec
	encoded := base64.StdEncoding.EncodeToString(marshaled)
	os.Setenv("GRAPH_CONFIG_BASE64", encoded)
	os.Setenv("JWT_SECRET_KEY", "test-secret-key-for-testing-purposes-only")
	os.Setenv("COOKIE_DOMAIN", "localhost")

	var mu sync.Mutex
	var capturedAdapter bool

	wrappedServe := DefaultServe
	DefaultServe = func(s *Server) error {
		mu.Lock()
		capturedAdapter = true
		mu.Unlock()
		return wrappedServe(s)
	}

	oldStartFunc := defaultStartFunc
	defaultStartFunc = func(handler any) {
		mu.Lock()
		capturedAdapter = true
		mu.Unlock()
	}

	err = Start()

	mu.Lock()
	called := capturedAdapter
	mu.Unlock()

	if !called {
		t.Error("Expected Lambda adapter to be invoked when AWS_LAMBDA_FUNCTION_NAME is set, but it was not triggered")
	}

	DefaultServe = wrappedServe
	defaultStartFunc = oldStartFunc

	if err != nil {
		t.Errorf("Start returned unexpected error: %v", err)
	}
}
