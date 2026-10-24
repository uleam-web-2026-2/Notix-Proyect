package repository

import (
	"errors"
	"fmt"

	"notix/internal/models"

	"gorm.io/gorm"
)

type AccessRepository struct {
	db *gorm.DB
}

func NewAccessRepository(db *gorm.DB) *AccessRepository {
	return &AccessRepository{db: db}
}

func (r *AccessRepository) RoleForBoard(userID, boardID uint) (string, error) {
	var board models.Tablero
	if err := r.db.First(&board, boardID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrNotFound
	} else if err != nil {
		return "", fmt.Errorf("find board %d for access check: %w", boardID, err)
	}
	if board.CreadorID == userID {
		return models.RolPropietario, nil
	}

	var share models.TableroCompartido
	err := r.db.Where("tablero_id = ? AND usuario_id = ? AND estado_invitacion = ?", boardID, userID, models.InvitacionActiva).
		First(&share).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("find board membership: %w", err)
	}
	return share.RolEnTablero, nil
}

func (r *AccessRepository) IsUserOnBoard(userID, boardID uint) (bool, error) {
	role, err := r.RoleForBoard(userID, boardID)
	if err != nil {
		return false, err
	}
	return role == models.RolPropietario || role == models.RolColaborador, nil
}
