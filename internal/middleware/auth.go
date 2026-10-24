package middleware

import (
	"context"
	"net/http"
	"strings"

	"notix/internal/services"
)

type principalKey struct{}

type Principal struct {
	UserID uint
	Role   string
}

type Authenticator interface {
	Authenticate(token string) (services.TokenClaims, error)
}

func RequireAuth(auth Authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
			writeUnauthorized(w)
			return
		}

		claims, err := auth.Authenticate(strings.TrimSpace(token))
		if err != nil {
			writeUnauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), principalKey{}, Principal{UserID: claims.UserID, Role: claims.Role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PrincipalFrom(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(Principal)
	return principal, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"Se requiere un token Bearer válido."}}`))
}
