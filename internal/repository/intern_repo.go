package repository

import (
	"errors"
	"time"

	"backend-dashboard/internal/domain"

	"gorm.io/gorm"
)

type internRepository struct {
	db *gorm.DB
}

// NewInternRepository creates a new intern repository
func NewInternRepository(db *gorm.DB) domain.InternRepository {
	return &internRepository{db: db}
}

// Create creates a new intern profile
func (r *internRepository) Create(userID, picID uint, batch, division, university, major string, startDate, endDate time.Time) (*domain.InternProfile, error) {
	profile := &domain.InternProfile{
		UserID:     userID,
		PICID:      picID,
		Batch:      batch,
		Division:   division,
		University: university,
		Major:      major,
		StartDate:  startDate,
		EndDate:    endDate,
		CreatedAt:  time.Now(),
	}

	if err := r.db.Create(profile).Error; err != nil {
		return nil, err
	}

	return profile, nil
}

// GetByID gets an intern profile by ID
func (r *internRepository) GetByID(id uint) (*domain.InternProfile, error) {
	var profile domain.InternProfile
	err := r.db.Preload("User").Preload("PIC").First(&profile, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("intern profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

// GetByUserID gets an intern profile by user ID
func (r *internRepository) GetByUserID(userID uint) (*domain.InternProfile, error) {
	var profile domain.InternProfile
	err := r.db.Preload("User").Preload("PIC").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("intern profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

// GetAll gets all intern profiles with pagination
func (r *internRepository) GetAll(page, limit int) ([]domain.InternProfile, int64, error) {
	var profiles []domain.InternProfile
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&domain.InternProfile{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := r.db.Preload("User").Preload("PIC").
		Offset(offset).
		Limit(limit).
		Find(&profiles).Error

	if err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}

// Update updates an intern profile
func (r *internRepository) Update(id uint, batch, division, university, major string) (*domain.InternProfile, error) {
	var profile domain.InternProfile
	if err := r.db.First(&profile, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("intern profile not found")
		}
		return nil, err
	}

	profile.Batch = batch
	profile.Division = division
	profile.University = university
	profile.Major = major

	if err := r.db.Save(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}
