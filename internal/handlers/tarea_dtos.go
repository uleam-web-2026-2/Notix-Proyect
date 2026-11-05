package handlers

import (
	"time"

	"notix/internal/models"
)

// CrearTareaRequest representa los únicos datos de tarea aceptados por POST /tareas.
type CrearTareaRequest struct {
	Titulo    string `json:"titulo"`
	Contenido string `json:"contenido"`
	ColumnaID uint   `json:"columna_id"`
}

func (p CrearTareaRequest) aEntidad() models.Tarea {
	return models.Tarea{
		Titulo:    p.Titulo,
		Contenido: p.Contenido,
		ColumnaID: p.ColumnaID,
	}
}

// TareaResponse define los campos de tarea expuestos por HTTP.
type TareaResponse struct {
	ID        uint      `json:"id"`
	Titulo    string    `json:"titulo"`
	Contenido string    `json:"contenido"`
	Estado    string    `json:"estado"`
	ColumnaID uint      `json:"columna_id"`
	CreadorID uint      `json:"creador_id"`
	Creado    time.Time `json:"creado"`
}

func nuevaRespuesta(t models.Tarea) TareaResponse {
	return TareaResponse{
		ID:        t.ID,
		Titulo:    t.Titulo,
		Contenido: t.Contenido,
		Estado:    t.Estado,
		ColumnaID: t.ColumnaID,
		CreadorID: t.CreadorID,
		Creado:    t.Creado,
	}
}

func listaRespuesta(tareas []models.Tarea) []TareaResponse {
	respuestas := make([]TareaResponse, 0, len(tareas))
	for _, tarea := range tareas {
		respuestas = append(respuestas, nuevaRespuesta(tarea))
	}
	return respuestas
}

// CambiarEstadoRequest contiene el estado solicitado por PATCH /tareas/{id}/estado.
type CambiarEstadoRequest struct {
	Estado string `json:"estado"`
}
