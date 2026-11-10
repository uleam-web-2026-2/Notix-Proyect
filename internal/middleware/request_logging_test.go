package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestLogRequestsRecordsMethodPathAndStatus(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.Handler
		method     string
		path       string
		wantStatus int
	}{
		{
			name: "explicit status",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			}),
			method:     http.MethodGet,
			path:       "/tableros/1/tareas",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "implicit status",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("ok"))
			}),
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			logger := log.New(&output, "", 0)
			handler := LogRequests(logger, tt.handler)
			request := httptest.NewRequest(tt.method, tt.path+"?token=not-logged", nil)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("response status = %d, want %d", response.Code, tt.wantStatus)
			}
			want := tt.method + " " + tt.path + " " + strconv.Itoa(tt.wantStatus)
			got := strings.TrimSpace(output.String())
			if !strings.HasSuffix(got, want) {
				t.Fatalf("log = %q, want it to end with %q", got, want)
			}
			if strings.Contains(got, "token=") {
				t.Fatalf("log contains query parameters: %q", got)
			}
		})
	}
}
