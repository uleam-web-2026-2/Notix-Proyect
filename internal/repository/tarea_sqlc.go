package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"notix/internal/models"
	"notix/internal/repository/filas"
)

// TareaSqlcRepository usa sqlc para buscar por ID y cambiar el estado. El resto
// de las operaciones se delega al repositorio GORM embebido.
type TareaSqlcRepository struct {
	*TareaRepository
	queries *filas.Queries
}

func NewTareaSqlcRepository(gormRepository *TareaRepository, db *sql.DB) *TareaSqlcRepository {
	return &TareaSqlcRepository{
		TareaRepository: gormRepository,
		queries:         filas.New(db),
	}
}

func (r *TareaSqlcRepository) FindByID(id uint) (models.Tarea, error) {
	row, err := r.queries.BuscarTarea(context.Background(), int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.Tarea{}, ErrNotFound
	}
	if err != nil {
		return models.Tarea{}, fmt.Errorf("sqlc find task %d: %w", id, err)
	}
	column, err := r.TareaRepository.FindColumn(uint(row.ColumnaID))
	if err != nil {
		return models.Tarea{}, err
	}
	task := models.Tarea{
		ID:        uint(row.ID),
		Titulo:    row.Titulo,
		Estado:    row.Estado,
		ColumnaID: uint(row.ColumnaID),
		CreadorID: uint(row.CreadorID),
		Columna:   column,
	}
	if row.Contenido.Valid {
		task.Contenido = row.Contenido.String
	}
	if row.Creado.Valid {
		task.Creado = row.Creado.Time
	}
	return task, nil
}

func (r *TareaSqlcRepository) UpdateState(id uint, from, to string) error {
	affected, err := r.queries.ActualizarEstadoTarea(context.Background(), filas.ActualizarEstadoTareaParams{
		ID:       int64(id),
		Estado:   to,
		Estado_2: from,
	})
	if err != nil {
		return fmt.Errorf("sqlc update task %d state: %w", id, err)
	}
	if affected > 0 {
		return nil
	}
	if _, err := r.queries.BuscarTarea(context.Background(), int64(id)); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("sqlc check task %d after state update: %w", id, err)
	}
	return ErrConflict
}
