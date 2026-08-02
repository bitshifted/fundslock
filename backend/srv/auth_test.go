package srv

import (
	"bitshifted/fundslock-be/model"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_createNonce(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/auth/nonce", nil)
	rec := httptest.NewRecorder()

	createNonce(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	var nonce model.Nonce
	err := json.Unmarshal(body, &nonce)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	if nonce.Nonce == "" {
		t.Errorf("Nonce is empty")
	}

}
