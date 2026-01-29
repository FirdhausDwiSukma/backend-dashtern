package repository

import (
	"backend-dashboard/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type postgresRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *postgresRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}
