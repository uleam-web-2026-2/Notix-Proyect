package services

import "notix/internal/models"

// TareaRepository contiene las operaciones de persistencia que usa TareaService.
type TareaRepository interface {
	ListByBoard(boardID uint) ([]models.Tarea, error)
	FindByID(id uint) (models.Tarea, error)
	Create(task *models.Tarea) error
	UpdateState(id uint, from, to string) error
	FindColumn(id uint) (models.Columna, error)
}

// BoardAccessRepository expone las comprobaciones de acceso usadas por TareaService.
type BoardAccessRepository interface {
	IsUserOnBoard(userID, boardID uint) (bool, error)
	RoleForBoard(userID, boardID uint) (string, error)
}
