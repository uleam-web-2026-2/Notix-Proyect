package repository

import (
	"errors"
	"fmt"

	"notix/internal/models"

	"gorm.io/gorm"
)

type TareaRepository struct {
	db *gorm.DB
}

func NewTareaRepository(db *gorm.DB) *TareaRepository {
	return &TareaRepository{db: db}
}

func (r *TareaRepository) ListByBoard(boardID uint) ([]models.Tarea, error) {
	var tasks []models.Tarea
	err := r.db.
		Joins("JOIN columnas ON columnas.id = tareas.columna_id").
		Where("columnas.tablero_id = ?", boardID).
		Preload("Columna").
		Preload("Creador").
		Order("columnas.orden ASC, tareas.creado ASC, tareas.id ASC").
		Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("list tasks for board %d: %w", boardID, err)
	}
	if tasks == nil {
		tasks = []models.Tarea{}
	}
	return tasks, nil
}

func (r *TareaRepository) FindByID(id uint) (models.Tarea, error) {
	var task models.Tarea
	err := r.db.Preload("Columna").Preload("Creador").First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Tarea{}, ErrNotFound
	}
	if err != nil {
		return models.Tarea{}, fmt.Errorf("find task %d: %w", id, err)
	}
	return task, nil
}

func (r *TareaRepository) Create(task *models.Tarea) error {
	if err := r.db.Create(task).Error; err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (r *TareaRepository) UpdateState(id uint, from, to string) error {
	result := r.db.Model(&models.Tarea{}).
		Where("id = ? AND estado = ?", id, from).
		Update("estado", to)
	if result.Error != nil {
		return fmt.Errorf("update task %d state: %w", id, result.Error)
	}
	if result.RowsAffected > 0 {
		return nil
	}

	var exists bool
	if err := r.db.Model(&models.Tarea{}).Select("count(*) > 0").Where("id = ?", id).Scan(&exists).Error; err != nil {
		return fmt.Errorf("check task %d after concurrent update: %w", id, err)
	}
	if !exists {
		return ErrNotFound
	}
	return ErrConflict
}

func (r *TareaRepository) FindColumn(id uint) (models.Columna, error) {
	var column models.Columna
	err := r.db.Preload("Tablero").First(&column, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Columna{}, ErrNotFound
	}
	if err != nil {
		return models.Columna{}, fmt.Errorf("find column %d: %w", id, err)
	}
	return column, nil
}
