package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainRoutes(t *testing.T) {
	mux := http.NewServeMux()
	MainRoutes(mux)

	// GET /
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET / returned %v, expected %v", rr.Code, http.StatusOK)
	}

	// POST /pull-request
	requestBody, err := json.Marshal(map[string]string{
		"action": "test",
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req, err = http.NewRequest("POST", "/pull-request", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("POST /pull-request returned %v, expected %v", rr.Code, http.StatusOK)
	}

}
