package httpapi

import "net/http"

func withFrontendCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != frontendOrigin {
			next.ServeHTTP(w, r)
			return
		}

		header := w.Header()
		header.Set("Access-Control-Allow-Origin", frontendOrigin)
		header.Add("Vary", "Origin")

		if r.Method == http.MethodOptions {
			header.Add("Vary", "Access-Control-Request-Method")
			header.Add("Vary", "Access-Control-Request-Headers")
			header.Set("Access-Control-Allow-Methods", "GET, POST, PATCH")
			header.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			header.Set("Access-Control-Max-Age", "300")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
