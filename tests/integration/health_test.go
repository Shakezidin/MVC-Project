//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthEndpoint is an integration test skeleton.
// Run with: go test -tags=integration ./tests/integration/...
func TestHealthEndpoint(t *testing.T) {
	baseURL := getBaseURL()

	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("server not running, skipping integration test: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status 200 or 503, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if _, ok := result["success"]; !ok {
		t.Error("response missing 'success' field")
	}
}

// TestLiveEndpoint verifies liveness probe.
func TestLiveEndpoint(t *testing.T) {
	baseURL := getBaseURL()

	resp, err := http.Get(baseURL + "/live")
	if err != nil {
		t.Skipf("server not running, skipping integration test: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func getBaseURL() string {
	return "http://localhost:8080"
}

// Example of in-process integration test using httptest
func TestLiveEndpointInProcess(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"message":"service is alive","data":{"status":"alive"}}`))
	})

	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
