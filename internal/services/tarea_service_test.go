package services

import (
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
