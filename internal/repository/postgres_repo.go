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
	result := r.db.Preload("Role").Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *postgresRepo) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	result := r.db.Preload("Role").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *postgresRepo) GetAll(page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	// Count total records
	r.db.Model(&domain.User{}).Count(&total)

	// Get paginated results
	offset := (page - 1) * limit
	result := r.db.Preload("Role").Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

func (r *postgresRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *postgresRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *postgresRepo) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
