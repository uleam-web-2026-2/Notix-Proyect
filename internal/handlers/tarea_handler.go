package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"notix/internal/middleware"
	"notix/internal/services"
)

type TareaHandler struct {
	tasks *services.TareaService
}

func NewTareaHandler(tasks *services.TareaService) *TareaHandler {
	return &TareaHandler{tasks: tasks}
}

func (h *TareaHandler) ListByBoard(w http.ResponseWriter, r *http.Request) {
	boardID, ok := parseID(w, r.PathValue("tableroID"), "tablero_id")
	if !ok {
		return
	}
	principal, ok := middleware.PrincipalFrom(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: errorDetail{Code: "unauthorized", Message: "Se requiere iniciar sesión."}})
		return
	}
	tasks, err := h.tasks.ListByBoard(principal.UserID, boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TareaHandler) Get(w http.ResponseWriter, r *http.Request) {
	taskID, ok := parseID(w, r.PathValue("id"), "id")
	if !ok {
		return
	}
	principal, ok := middleware.PrincipalFrom(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: errorDetail{Code: "unauthorized", Message: "Se requiere iniciar sesión."}})
		return
	}
	task, err := h.tasks.Find(principal.UserID, taskID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TareaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Titulo    string `json:"titulo"`
		Contenido string `json:"contenido"`
		ColumnaID uint   `json:"columna_id"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	principal, ok := middleware.PrincipalFrom(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: errorDetail{Code: "unauthorized", Message: "Se requiere iniciar sesión."}})
		return
	}
	task, err := h.tasks.Create(principal.UserID, input.ColumnaID, input.Titulo, input.Contenido)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TareaHandler) ChangeState(w http.ResponseWriter, r *http.Request) {
	taskID, ok := parseID(w, r.PathValue("id"), "id")
	if !ok {
		return
	}
	var input struct {
		Estado string `json:"estado"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	principal, ok := middleware.PrincipalFrom(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: errorDetail{Code: "unauthorized", Message: "Se requiere iniciar sesión."}})
		return
	}
	task, err := h.tasks.ChangeState(principal.UserID, taskID, strings.TrimSpace(input.Estado))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func parseID(w http.ResponseWriter, raw, name string) (uint, bool) {
	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: errorDetail{Code: "invalid_id", Message: name + " debe ser un entero positivo."}})
		return 0, false
	}
	return uint(id), true
}
