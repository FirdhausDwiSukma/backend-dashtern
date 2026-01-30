package usecase

import (
	"time"

	"backend-dashboard/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type internUsecase struct {
	internRepo domain.InternRepository
	userRepo   domain.UserRepository
}

// NewInternUsecase creates a new intern usecase
func NewInternUsecase(internRepo domain.InternRepository, userRepo domain.UserRepository) domain.InternUsecase {
	return &internUsecase{
		internRepo: internRepo,
		userRepo:   userRepo,
	}
}

// CreateIntern creates a new intern user with profile
func (u *internUsecase) CreateIntern(fullName, username, email, password string, picID uint, batch, division, university, major string, startDate, endDate time.Time) (*domain.User, *domain.InternProfile, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	// Get intern role ID (assuming role ID 4 is for intern based on seeding)
	// In production, you should fetch this from roles table
	internRoleID := uint(4)

	// Create user account
	user := &domain.User{
		FullName:     fullName,
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       internRoleID,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user
	err = u.userRepo.Create(user)
	if err != nil {
		return nil, nil, err
	}

	// Create intern profile
	profile, err := u.internRepo.Create(user.ID, picID, batch, division, university, major, startDate, endDate)
	if err != nil {
		// If profile creation fails, we should ideally rollback user creation
		// For now, just return the error
		return nil, nil, err
	}

	return user, profile, nil
}

// GetInternByID gets an intern profile by ID
func (u *internUsecase) GetInternByID(id uint) (*domain.InternProfile, error) {
	profile, err := u.internRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// GetAllInterns gets all interns with pagination
func (u *internUsecase) GetAllInterns(page, limit int) ([]domain.InternProfile, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	profiles, total, err := u.internRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}
