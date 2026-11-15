package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	//mock request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// run handler
	healthHandler(w, req)

	// get result
	res := w.Result()
	defer res.Body.Close()

	// test "ok" status
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	// test content type equal to application/json
	if contentType := res.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	// test response body is json
	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// test response body status is ok
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %s", body["status"])
	}

}
