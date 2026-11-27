package transporte

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"notix/internal/models"
	"notix/internal/services"
)

type authenticadorDePrueba struct {
	tokens map[string]services.TokenClaims
}

func (a authenticadorDePrueba) Authenticate(token string) (services.TokenClaims, error) {
	claims, ok := a.tokens[token]
	if !ok {
		return services.TokenClaims{}, errors.New("token inválido")
	}
	return claims, nil
}

func TestWebSocketSinTokenResponde401(t *testing.T) {
	server, _ := newTestServer(t)
	response, err := server.Client().Get(server.URL + "/ws")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}
}

func TestWebSocketOrigenNoPermitidoResponde403(t *testing.T) {
	server, _ := newTestServer(t)
	_, response, err := websocket.Dial(context.Background(), websocketURL(server.URL)+"?token=ana", &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://otro-origen.test"}},
	})
	if err == nil {
		t.Fatal("Dial() succeeded for an origin that should be rejected")
	}
	if response == nil || response.StatusCode != http.StatusForbidden {
		if response == nil {
			t.Fatal("rejected handshake did not return an HTTP response")
		}
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusForbidden)
	}
}

func TestWebSocketAvisoSoloLlegaAlCreador(t *testing.T) {
	server, avisos := newTestServer(t)
	ana := dialTestClient(t, server.URL, "ana")
	defer ana.Close(websocket.StatusNormalClosure, "")
	eva := dialTestClient(t, server.URL, "eva")
	defer eva.Close(websocket.StatusNormalClosure, "")

	waitForClients(t, avisos, "1", "2")
	avisos.Enviar("1", services.Aviso{Tipo: "estado_cambiado", ID: 23, Estado: "en_proceso"})

	readCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	tipo, contenido, err := ana.Read(readCtx)
	if err != nil {
		t.Fatalf("creator did not receive the notification: %v", err)
	}
	if tipo != websocket.MessageText || !strings.Contains(string(contenido), `"id":23`) {
		t.Fatalf("unexpected notification: type=%v body=%s", tipo, contenido)
	}

	noReadCtx, cancelNoRead := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelNoRead()
	if _, _, err := eva.Read(noReadCtx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("second collaborator read error = %v, want timeout with no notification", err)
	}
}

func TestRetoAvisaATodosLosPropietariosConectados(t *testing.T) {
	server, avisos := newTestServer(t)
	creador := dialTestClient(t, server.URL, "ana")
	defer creador.Close(websocket.StatusNormalClosure, "")
	propietario1 := dialTestClient(t, server.URL, "luis")
	defer propietario1.Close(websocket.StatusNormalClosure, "")
	propietario2 := dialTestClient(t, server.URL, "marta")
	defer propietario2.Close(websocket.StatusNormalClosure, "")

	waitForClients(t, avisos, "1", "3", "4")
	avisos.AvisarCreacion(1, 23)

	want := services.Aviso{Tipo: "tarea_nueva", ID: 23}
	assertReceivesAviso(t, propietario1, want)
	assertReceivesAviso(t, propietario2, want)
	assertReceivesNoAviso(t, creador)
}

func TestRetoExcluyeAlPropietarioCreador(t *testing.T) {
	server, avisos := newTestServer(t)
	creador := dialTestClient(t, server.URL, "propietario-creador")
	defer creador.Close(websocket.StatusNormalClosure, "")
	propietario := dialTestClient(t, server.URL, "luis")
	defer propietario.Close(websocket.StatusNormalClosure, "")
	colaborador := dialTestClient(t, server.URL, "ana")
	defer colaborador.Close(websocket.StatusNormalClosure, "")

	waitForClients(t, avisos, "1", "3", "5")
	avisos.AvisarCreacion(5, 24)

	assertReceivesAviso(t, propietario, services.Aviso{Tipo: "tarea_nueva", ID: 24})
	assertReceivesNoAviso(t, creador)
	assertReceivesNoAviso(t, colaborador)
}

func assertReceivesAviso(t *testing.T, conn *websocket.Conn, want services.Aviso) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	messageType, contents, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("client did not receive notification: %v", err)
	}
	if messageType != websocket.MessageText {
		t.Fatalf("message type = %v, want text", messageType)
	}
	var got services.Aviso
	if err := json.Unmarshal(contents, &got); err != nil {
		t.Fatalf("unmarshal notification: %v", err)
	}
	if got != want {
		t.Fatalf("notification = %+v, want %+v", got, want)
	}
}

func assertReceivesNoAviso(t *testing.T, conn *websocket.Conn) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if _, _, err := conn.Read(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("client read error = %v, want no notification before timeout", err)
	}
}

func newTestServer(t *testing.T) (*httptest.Server, *Avisos) {
	t.Helper()
	avisos := NuevosAvisos(authenticadorDePrueba{tokens: map[string]services.TokenClaims{
		"ana":                 {UserID: 1, Role: models.RolColaborador},
		"eva":                 {UserID: 2, Role: models.RolColaborador},
		"luis":                {UserID: 3, Role: models.RolPropietario},
		"marta":               {UserID: 4, Role: models.RolPropietario},
		"propietario-creador": {UserID: 5, Role: models.RolPropietario},
	}}, "http://localhost:5173")
	server := httptest.NewServer(http.HandlerFunc(avisos.Conectar))
	t.Cleanup(server.Close)
	return server, avisos
}

func dialTestClient(t *testing.T, serverURL, token string) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.Dial(context.Background(), websocketURL(serverURL)+"?token="+token, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://localhost:5173"}},
	})
	if err != nil {
		t.Fatalf("dial websocket as %s: %v", token, err)
	}
	return conn
}

func waitForClients(t *testing.T, avisos *Avisos, users ...string) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		avisos.mu.Lock()
		ready := true
		for _, user := range users {
			if len(avisos.canales[user]) == 0 {
				ready = false
				break
			}
		}
		avisos.mu.Unlock()
		if ready {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("websocket clients did not connect in time")
}

func websocketURL(serverURL string) string {
	return "ws" + strings.TrimPrefix(serverURL, "http") + "/ws"
}
