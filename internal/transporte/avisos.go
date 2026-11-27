package transporte

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/coder/websocket"

	"notix/internal/models"
	"notix/internal/services"
)

type Autenticador interface {
	Authenticate(token string) (services.TokenClaims, error)
}

type Avisos struct {
	autenticador  Autenticador
	origen        string
	mu            sync.Mutex
	canales       map[string]map[*canalCliente]struct{}
	canalesPorRol map[string]map[*canalCliente]struct{}
}

type canalCliente struct {
	conn   *websocket.Conn
	userID uint
	avisos chan services.Aviso
}

func NuevosAvisos(autenticador Autenticador, origen string) *Avisos {
	return &Avisos{
		autenticador:  autenticador,
		origen:        origen,
		canales:       make(map[string]map[*canalCliente]struct{}),
		canalesPorRol: make(map[string]map[*canalCliente]struct{}),
	}
}

func (a *Avisos) Enviar(quien string, aviso services.Aviso) {
	a.enviar(a.clientesPorUsuario(quien), aviso)
}

func (a *Avisos) AvisarCreacion(creadorID, registroID uint) {
	clientes := a.clientesPorRol(models.RolPropietario, creadorID)
	a.enviar(clientes, services.Aviso{Tipo: "tarea_nueva", ID: registroID})
}

func (a *Avisos) clientesPorUsuario(quien string) []*canalCliente {
	a.mu.Lock()
	defer a.mu.Unlock()
	clientes := make([]*canalCliente, 0, len(a.canales[quien]))
	for cliente := range a.canales[quien] {
		clientes = append(clientes, cliente)
	}
	return clientes
}

func (a *Avisos) clientesPorRol(rol string, excluirID uint) []*canalCliente {
	a.mu.Lock()
	defer a.mu.Unlock()
	clientes := make([]*canalCliente, 0, len(a.canalesPorRol[rol]))
	for cliente := range a.canalesPorRol[rol] {
		if cliente.userID != excluirID {
			clientes = append(clientes, cliente)
		}
	}
	return clientes
}

func (a *Avisos) enviar(clientes []*canalCliente, aviso services.Aviso) {
	for _, cliente := range clientes {
		select {
		case cliente.avisos <- aviso:
		default:
			// Una conexión lenta no detiene el cambio de estado.
		}
	}
}

func (a *Avisos) Conectar(w http.ResponseWriter, r *http.Request) {
	claims, err := a.autenticador.Authenticate(r.URL.Query().Get("token"))
	if err != nil || claims.UserID == 0 {
		writeUnauthorized(w)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{a.origen},
	})
	if err != nil {
		return
	}
	quien := strconv.FormatUint(uint64(claims.UserID), 10)
	cliente := &canalCliente{conn: conn, userID: claims.UserID, avisos: make(chan services.Aviso, 16)}
	readContext := conn.CloseRead(context.Background())
	a.agregar(quien, claims.Role, cliente)
	go a.escribir(readContext, cliente)
	defer func() {
		a.quitar(quien, claims.Role, cliente)
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}()
	<-readContext.Done()
}

func (a *Avisos) escribir(ctx context.Context, cliente *canalCliente) {
	for {
		select {
		case <-ctx.Done():
			return
		case aviso := <-cliente.avisos:
			contenido, err := json.Marshal(aviso)
			if err != nil {
				continue
			}
			writeContext, cancel := context.WithTimeout(ctx, 5*time.Second)
			err = cliente.conn.Write(writeContext, websocket.MessageText, contenido)
			cancel()
			if err != nil {
				_ = cliente.conn.Close(websocket.StatusInternalError, "")
				return
			}
		}
	}
}

func (a *Avisos) agregar(quien, rol string, cliente *canalCliente) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.canales[quien] == nil {
		a.canales[quien] = make(map[*canalCliente]struct{})
	}
	a.canales[quien][cliente] = struct{}{}
	if a.canalesPorRol[rol] == nil {
		a.canalesPorRol[rol] = make(map[*canalCliente]struct{})
	}
	a.canalesPorRol[rol][cliente] = struct{}{}
}

func (a *Avisos) quitar(quien, rol string, cliente *canalCliente) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.canales[quien], cliente)
	if len(a.canales[quien]) == 0 {
		delete(a.canales, quien)
	}
	delete(a.canalesPorRol[rol], cliente)
	if len(a.canalesPorRol[rol]) == 0 {
		delete(a.canalesPorRol, rol)
	}
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"Se requiere un token Bearer válido."}}`))
}
