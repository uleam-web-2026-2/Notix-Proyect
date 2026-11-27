package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFrontendCORSAllowsPagePreflight(t *testing.T) {
	handler := withFrontendCORS(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("preflight should not reach the API handler")
	}))
	request := httptest.NewRequest(http.MethodOptions, "/login", nil)
	request.Header.Set("Origin", frontendOrigin)
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	request.Header.Set("Access-Control-Request-Headers", "content-type")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != frontendOrigin {
		t.Fatalf("allow-origin = %q, want %q", got, frontendOrigin)
	}
	if got := response.Header().Get("Access-Control-Allow-Methods"); got != "GET, POST, PATCH" {
		t.Fatalf("allow-methods = %q", got)
	}
	if got := response.Header().Get("Access-Control-Allow-Headers"); got != "Authorization, Content-Type" {
		t.Fatalf("allow-headers = %q", got)
	}
}

func TestFrontendCORSDoesNotAllowOtherOrigins(t *testing.T) {
	handler := withFrontendCORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	request := httptest.NewRequest(http.MethodOptions, "/login", nil)
	request.Header.Set("Origin", "http://otro-origen.test")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("allow-origin = %q, want no CORS permission", got)
	}
}
