package repository

import (
	"errors"
	"fmt"

	"notix/internal/models"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) FindByEmail(email string) (models.Cuenta, error) {
	var account models.Cuenta
	err := r.db.Preload("Usuario").Where("correo = ?", email).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Cuenta{}, ErrNotFound
	}
	if err != nil {
		return models.Cuenta{}, fmt.Errorf("find login account: %w", err)
	}
	return account, nil
}
