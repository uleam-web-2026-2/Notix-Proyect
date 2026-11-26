package services

import (
	"errors"
	"testing"

	"notix/internal/models"
)

func TestValidTransition(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want bool
	}{
		{name: "start work", from: models.EstadoPendiente, to: models.EstadoEnProceso, want: true},
		{name: "pause work", from: models.EstadoEnProceso, to: models.EstadoPendiente, want: true},
		{name: "finish work", from: models.EstadoEnProceso, to: models.EstadoCompletado, want: true},
		{name: "reopen work", from: models.EstadoCompletado, to: models.EstadoEnProceso, want: true},
		{name: "cannot complete without starting", from: models.EstadoPendiente, to: models.EstadoCompletado, want: false},
		{name: "cannot repeat current state", from: models.EstadoPendiente, to: models.EstadoPendiente, want: false},
		{name: "unknown state", from: "desconocido", to: models.EstadoPendiente, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := validTransition(test.from, test.to); got != test.want {
				t.Fatalf("validTransition(%q, %q) = %t, want %t", test.from, test.to, got, test.want)
			}
		})
	}
}

type fakeTareaStore struct {
	task    models.Tarea
	updated bool
}

func (f *fakeTareaStore) ListByBoard(uint) ([]models.Tarea, error) { return nil, nil }
func (f *fakeTareaStore) FindByID(uint) (models.Tarea, error)      { return f.task, nil }
func (f *fakeTareaStore) FindColumn(uint) (models.Columna, error)  { return models.Columna{}, nil }
func (f *fakeTareaStore) Create(*models.Tarea) error               { return nil }
func (f *fakeTareaStore) UpdateState(_ uint, _, to string) error {
	f.task.Estado = to
	f.updated = true
	return nil
}

type fakeAccessStore struct{}

func (fakeAccessStore) IsUserOnBoard(uint, uint) (bool, error) { return true, nil }
func (fakeAccessStore) RoleForBoard(uint, uint) (string, error) {
	return models.RolPropietario, nil
}

type fakeAvisador struct {
	who    string
	aviso  Aviso
	called bool
	store  *fakeTareaStore
}

func (f *fakeAvisador) Enviar(who string, aviso Aviso) {
	f.called = true
	f.who = who
	f.aviso = aviso
	if !f.store.updated {
		panic("notification was sent before the state was saved")
	}
}

func TestCambioEstadoAvisaAlCreadorDespuesDeGuardar(t *testing.T) {
	store := &fakeTareaStore{task: models.Tarea{
		ID: 23, Estado: models.EstadoPendiente, CreadorID: 7,
		Columna: models.Columna{TableroID: 4},
	}}
	avisador := &fakeAvisador{store: store}
	service := NewTareaService(store, fakeAccessStore{})
	service.AvisarCon(avisador)

	task, err := service.ChangeState(9, 23, models.EstadoEnProceso)
	if err != nil {
		t.Fatalf("ChangeState() error = %v", err)
	}
	if task.Estado != models.EstadoEnProceso {
		t.Fatalf("task state = %q, want %q", task.Estado, models.EstadoEnProceso)
	}
	if !avisador.called || avisador.who != "7" || avisador.aviso != (Aviso{Tipo: "estado_cambiado", ID: 23, Estado: models.EstadoEnProceso}) {
		t.Fatalf("unexpected notification: called=%t who=%q aviso=%+v", avisador.called, avisador.who, avisador.aviso)
	}
}

func TestTransicionProhibidaNoAvisa(t *testing.T) {
	store := &fakeTareaStore{task: models.Tarea{
		ID: 23, Estado: models.EstadoPendiente, CreadorID: 7,
		Columna: models.Columna{TableroID: 4},
	}}
	avisador := &fakeAvisador{store: store}
	service := NewTareaService(store, fakeAccessStore{})
	service.AvisarCon(avisador)

	_, err := service.ChangeState(9, 23, models.EstadoCompletado)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("ChangeState() error = %v, want %v", err, ErrInvalidTransition)
	}
	if avisador.called {
		t.Fatal("prohibited transition sent a notification")
	}
}
