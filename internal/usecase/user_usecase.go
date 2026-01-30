package usecase

import (
	"backend-dashboard/internal/domain"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
	loginUC   *loginUsecase
}

func NewUserUsecase(userRepo domain.UserRepository, jwtSecret string) domain.UserUsecase {
	loginUC := NewLoginUsecase(userRepo, jwtSecret)
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		loginUC:   loginUC,
	}
}

func (u *userUsecase) Login(username, password string) (string, *domain.User, error) {
	return u.loginUC.Login(username, password)
}

func (u *userUsecase) GetAll(page, limit int) ([]domain.User, int64, error) {
	return u.userRepo.GetAll(page, limit)
}

func (u *userUsecase) GetByID(id uint) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *userUsecase) Create(fullName, username, email, password string, roleID uint) (*domain.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		FullName:     fullName,
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       roleID,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Reload with role
	return u.userRepo.GetByID(user.ID)
}

func (u *userUsecase) Update(id uint, fullName, email string, status string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.FullName = fullName
	user.Email = email
	user.Status = status
	user.UpdatedAt = time.Now()

	err = u.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) Deactivate(id uint) error {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	user.Status = "inactive"
	user.UpdatedAt = time.Now()

	return u.userRepo.Update(user)
}

// HardDelete permanently deletes a user from the database
// WARNING: This is a destructive operation and should only be used by super_admin
// It will cascade delete all related data
func (u *userUsecase) HardDelete(id uint) error {
	// Check if user exists first
	_, err := u.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Permanently delete from database
	return u.userRepo.Delete(id)
}
