package httpapi

import (
	"net/http"

	"notix/config"
	"notix/internal/handlers"
	"notix/internal/middleware"
	"notix/internal/repository"
	"notix/internal/services"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg config.Config) http.Handler {
	tareaRepository := repository.NewTareaRepository(db)
	accessRepository := repository.NewAccessRepository(db)
	accountRepository := repository.NewAccountRepository(db)
	tareaService := services.NewTareaService(tareaRepository, accessRepository)
	authService := services.NewAuthService(accountRepository, cfg.JWTSecret, cfg.TokenExpiry)
	tareaHandler := handlers.NewTareaHandler(tareaService)
	authHandler := handlers.NewAuthHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle("GET /tableros/{tableroID}/tareas", middleware.RequireAuth(authService, http.HandlerFunc(tareaHandler.ListByBoard)))
	mux.Handle("GET /tareas/{id}", middleware.RequireAuth(authService, http.HandlerFunc(tareaHandler.Get)))
	mux.Handle("POST /tareas", middleware.RequireAuth(authService, http.HandlerFunc(tareaHandler.Create)))
	mux.Handle("PATCH /tareas/{id}/estado", middleware.RequireAuth(authService, http.HandlerFunc(tareaHandler.ChangeState)))
	return mux
}
