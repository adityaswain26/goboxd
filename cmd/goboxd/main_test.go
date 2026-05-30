package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	rr := httptest.NewRecorder()

	healthz(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestRunUnsupportedLanguage(t *testing.T) {
	body := `{
		"language":"brainfuck",
		"source":"++++",
		"expected_output":""
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/run",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	runHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}
func TestRunBadJSON(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/run",
		strings.NewReader("{invalid json"),
	)

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	runHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}
func TestReadyz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)

	rr := httptest.NewRecorder()

	readyz(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}
