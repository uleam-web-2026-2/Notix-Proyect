package handlers

import (
	"errors"
	"net/http"
	"strings"

	"notix/internal/services"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Correo string `json:"correo"`
		Clave  string `json:"clave"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	result, err := h.auth.Login(strings.TrimSpace(input.Correo), input.Clave)
	if errors.Is(err, services.ErrInvalidCredentials) {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: errorDetail{Code: "invalid_credentials", Message: "Correo o clave incorrectos."}})
		return
	}
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token":      result.Token,
		"token_type": "Bearer",
		"expires_at": result.ExpiresAt,
		"usuario": map[string]any{
			"id":     result.User.ID,
			"nombre": result.User.Nombre,
			"correo": result.User.Correo,
			"rol":    result.Role,
		},
	})
}
