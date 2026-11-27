package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogRequestsIncludesQueryAndResponseStatus(t *testing.T) {
	var output bytes.Buffer
	logger := log.New(&output, "", 0)
	handler := LogRequests(logger, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusSwitchingProtocols)
	}))
	request := httptest.NewRequest(http.MethodGet, "/ws?token=TEST", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if got, want := output.String(), "GET /ws?token=TEST 101\n"; got != want {
		t.Fatalf("log = %q, want %q", got, want)
	}
}

func TestLogRequestsUsesOKForImplicitStatus(t *testing.T) {
	var output bytes.Buffer
	logger := log.New(&output, "", 0)
	handler := LogRequests(logger, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/health", nil))

	if got, want := output.String(), "GET /health 200\n"; got != want {
		t.Fatalf("log = %q, want %q", got, want)
	}
}
