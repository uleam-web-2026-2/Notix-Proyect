package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"notix/internal/repository"
	"notix/internal/services"
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode JSON response: %v", err)
	}
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		writeJSON(w, http.StatusUnprocessableEntity, errorBody{Error: errorDetail{Code: "validation_error", Message: err.Error()}})
	case errors.Is(err, services.ErrInvalidTransition):
		writeJSON(w, http.StatusUnprocessableEntity, errorBody{Error: errorDetail{Code: "invalid_transition", Message: "La transición solicitada no está permitida."}})
	case errors.Is(err, services.ErrForbidden):
		writeJSON(w, http.StatusForbidden, errorBody{Error: errorDetail{Code: "forbidden", Message: "No tienes permiso para realizar esta acción."}})
	case errors.Is(err, repository.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{Error: errorDetail{Code: "not_found", Message: "No se encontró el recurso solicitado."}})
	case errors.Is(err, repository.ErrConflict):
		writeJSON(w, http.StatusConflict, errorBody{Error: errorDetail{Code: "conflict", Message: "La tarea cambió concurrentemente; vuelve a cargarla."}})
	default:
		log.Printf("request failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorBody{Error: errorDetail{Code: "internal_error", Message: "Ocurrió un error interno."}})
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: errorDetail{Code: "invalid_json", Message: "El cuerpo JSON no es válido."}})
		return false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: errorDetail{Code: "invalid_json", Message: "Envía un solo objeto JSON."}})
		return false
	}
	return true
}
