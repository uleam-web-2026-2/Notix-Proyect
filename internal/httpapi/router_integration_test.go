package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"notix/config"
	"notix/database"
	"notix/internal/models"
)

func TestTaskAPIIntegration(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil && !os.IsNotExist(err) {
		t.Fatalf("load root .env: %v", err)
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	db, err := database.Connect(cfg)
	if err != nil {
		if os.Getenv("RUN_DB_INTEGRATION") == "true" {
			t.Fatalf("PostgreSQL integration database is required: %v", err)
		}
		t.Skipf("PostgreSQL is unavailable; start Docker Compose to run this integration test: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get integration database connection pool: %v", err)
	}
	defer sqlDB.Close()
	if err := db.AutoMigrate(
		&models.Usuario{},
		&models.Tablero{},
		&models.Columna{},
		&models.Tarea{},
		&models.TableroCompartido{},
		&models.Comentario{},
		&models.Cuenta{},
	); err != nil {
		t.Fatalf("migrate integration tables: %v", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin integration transaction: %v", tx.Error)
	}
	defer tx.Rollback()

	owner, collaborator, board, column, tasks := createIntegrationFixture(t, tx)
	handler := NewRouter(tx, cfg)

	response := performRequest(handler, http.MethodGet, fmt.Sprintf("/tableros/%d/tareas", board.ID), nil, "")
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated list status = %d, want %d", response.Code, http.StatusUnauthorized)
	}

	collaboratorToken := loginForTest(t, handler, collaborator.Correo, "collaborator-password")
	ownerToken := loginForTest(t, handler, owner.Correo, "owner-password")

	response = performRequest(handler, http.MethodGet, fmt.Sprintf("/tableros/%d/tareas", board.ID), nil, collaboratorToken)
	if response.Code != http.StatusOK {
		t.Fatalf("collaborator list status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}

	createBody, _ := json.Marshal(map[string]any{
		"titulo":     "Tarea creada por integración",
		"contenido":  "La tarea empieza en pendiente.",
		"columna_id": column.ID,
	})
	response = performRequest(handler, http.MethodPost, "/tareas", createBody, collaboratorToken)
	if response.Code != http.StatusCreated {
		t.Fatalf("create task status = %d, want %d: %s", response.Code, http.StatusCreated, response.Body.String())
	}
	var created models.Tarea
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created task: %v", err)
	}
	if created.Estado != models.EstadoPendiente || created.CreadorID != collaborator.ID {
		t.Fatalf("created task has state %q and creator %d", created.Estado, created.CreadorID)
	}

	response = performRequest(handler, http.MethodGet, fmt.Sprintf("/tareas/%d", created.ID), nil, collaboratorToken)
	if response.Code != http.StatusOK {
		t.Fatalf("get task status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}

	stateBody, _ := json.Marshal(map[string]string{"estado": models.EstadoEnProceso})
	response = performRequest(handler, http.MethodPatch, fmt.Sprintf("/tareas/%d/estado", tasks[0].ID), stateBody, collaboratorToken)
	if response.Code != http.StatusForbidden {
		t.Fatalf("collaborator state change status = %d, want %d", response.Code, http.StatusForbidden)
	}

	response = performRequest(handler, http.MethodPatch, fmt.Sprintf("/tareas/%d/estado", tasks[0].ID), stateBody, ownerToken)
	if response.Code != http.StatusOK {
		t.Fatalf("owner state change status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}

	prohibitedBody, _ := json.Marshal(map[string]string{"estado": models.EstadoCompletado})
	response = performRequest(handler, http.MethodPatch, fmt.Sprintf("/tareas/%d/estado", tasks[1].ID), prohibitedBody, ownerToken)
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("prohibited state change status = %d, want %d", response.Code, http.StatusUnprocessableEntity)
	}
}

func createIntegrationFixture(t *testing.T, tx *gorm.DB) (models.Usuario, models.Usuario, models.Tablero, models.Columna, []models.Tarea) {
	t.Helper()
	suffix := time.Now().UnixNano()
	ownerEmail := fmt.Sprintf("owner-%d@integration.notix.test", suffix)
	collaboratorEmail := fmt.Sprintf("collaborator-%d@integration.notix.test", suffix)
	owner := createTestUser(t, tx, "Owner", ownerEmail, "owner-password", models.RolPropietario)
	collaborator := createTestUser(t, tx, "Collaborator", collaboratorEmail, "collaborator-password", models.RolColaborador)

	board := models.Tablero{Titulo: "Integration board", CreadorID: owner.ID}
	if err := tx.Create(&board).Error; err != nil {
		t.Fatalf("create integration board: %v", err)
	}
	column := models.Columna{Nombre: "Por hacer", Orden: 1, TableroID: board.ID}
	if err := tx.Create(&column).Error; err != nil {
		t.Fatalf("create integration column: %v", err)
	}
	share := models.TableroCompartido{
		TableroID:        board.ID,
		UsuarioID:        collaborator.ID,
		RolEnTablero:     models.RolColaborador,
		EstadoInvitacion: models.InvitacionActiva,
	}
	if err := tx.Create(&share).Error; err != nil {
		t.Fatalf("share integration board: %v", err)
	}
	tasks := []models.Tarea{
		{Titulo: "Valid transition", Estado: models.EstadoPendiente, ColumnaID: column.ID, CreadorID: collaborator.ID},
		{Titulo: "Prohibited transition", Estado: models.EstadoPendiente, ColumnaID: column.ID, CreadorID: owner.ID},
	}
	if err := tx.Create(&tasks).Error; err != nil {
		t.Fatalf("create integration tasks: %v", err)
	}
	return owner, collaborator, board, column, tasks
}

func createTestUser(t *testing.T, tx *gorm.DB, name, email, password, role string) models.Usuario {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash integration password: %v", err)
	}
	user := models.Usuario{Nombre: name, Correo: email, Clave: string(hash)}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("create integration user: %v", err)
	}
	account := models.Cuenta{UsuarioID: user.ID, Correo: email, ClaveHash: string(hash), Rol: role}
	if err := tx.Create(&account).Error; err != nil {
		t.Fatalf("create integration account: %v", err)
	}
	return user
}

func loginForTest(t *testing.T, handler http.Handler, email, password string) string {
	t.Helper()
	body, err := json.Marshal(map[string]string{"correo": email, "clave": password})
	if err != nil {
		t.Fatal(err)
	}
	response := performRequest(handler, http.MethodPost, "/login", body, "")
	if response.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if result.Token == "" {
		t.Fatal("login returned an empty token")
	}
	return result.Token
}

func performRequest(handler http.Handler, method, path string, body []byte, token string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
