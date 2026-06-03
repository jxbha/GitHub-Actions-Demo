package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()
	health(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("got %d, want %d", resp.Code, 200)

	}

}
