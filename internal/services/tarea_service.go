package services

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"notix/internal/models"
	"notix/internal/repository"
)

var (
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidTransition = errors.New("invalid state transition")
)

type TareaService struct {
	tasks  *repository.TareaRepository
	access *repository.AccessRepository
}

func NewTareaService(tasks *repository.TareaRepository, access *repository.AccessRepository) *TareaService {
	return &TareaService{tasks: tasks, access: access}
}

func (s *TareaService) ListByBoard(userID, boardID uint) ([]models.Tarea, error) {
	allowed, err := s.access.IsUserOnBoard(userID, boardID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrForbidden
	}
	return s.tasks.ListByBoard(boardID)
}

func (s *TareaService) Find(userID, taskID uint) (models.Tarea, error) {
	task, err := s.tasks.FindByID(taskID)
	if err != nil {
		return models.Tarea{}, err
	}
	allowed, err := s.access.IsUserOnBoard(userID, task.Columna.TableroID)
	if err != nil {
		return models.Tarea{}, err
	}
	if !allowed {
		return models.Tarea{}, ErrForbidden
	}
	return task, nil
}

func (s *TareaService) Create(userID, columnID uint, title, content string) (models.Tarea, error) {
	title = strings.TrimSpace(title)
	if title == "" || utf8.RuneCountInString(title) > 150 {
		return models.Tarea{}, fmt.Errorf("%w: titulo is required and must be at most 150 characters", ErrInvalidInput)
	}
	if columnID == 0 {
		return models.Tarea{}, fmt.Errorf("%w: columna_id must be a positive integer", ErrInvalidInput)
	}
	if utf8.RuneCountInString(content) > 5000 {
		return models.Tarea{}, fmt.Errorf("%w: contenido must be at most 5000 characters", ErrInvalidInput)
	}

	column, err := s.tasks.FindColumn(columnID)
	if err != nil {
		return models.Tarea{}, err
	}
	allowed, err := s.access.IsUserOnBoard(userID, column.TableroID)
	if err != nil {
		return models.Tarea{}, err
	}
	if !allowed {
		return models.Tarea{}, ErrForbidden
	}

	task := models.Tarea{
		Titulo:    title,
		Contenido: content,
		Estado:    models.EstadoPendiente,
		ColumnaID: columnID,
		CreadorID: userID,
	}
	if err := s.tasks.Create(&task); err != nil {
		return models.Tarea{}, err
	}
	return s.tasks.FindByID(task.ID)
}

func (s *TareaService) ChangeState(userID, taskID uint, target string) (models.Tarea, error) {
	task, err := s.tasks.FindByID(taskID)
	if err != nil {
		return models.Tarea{}, err
	}
	role, err := s.access.RoleForBoard(userID, task.Columna.TableroID)
	if err != nil {
		return models.Tarea{}, err
	}
	if role != models.RolPropietario {
		return models.Tarea{}, ErrForbidden
	}
	if !validTransition(task.Estado, target) {
		return models.Tarea{}, fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, task.Estado, target)
	}
	if err := s.tasks.UpdateState(task.ID, task.Estado, target); err != nil {
		return models.Tarea{}, err
	}
	return s.tasks.FindByID(task.ID)
}

func validTransition(from, to string) bool {
	switch from {
	case models.EstadoPendiente:
		return to == models.EstadoEnProceso
	case models.EstadoEnProceso:
		return to == models.EstadoPendiente || to == models.EstadoCompletado
	case models.EstadoCompletado:
		return to == models.EstadoEnProceso
	default:
		return false
	}
}
