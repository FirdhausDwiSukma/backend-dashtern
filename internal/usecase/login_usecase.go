package usecase

import (
	"backend-dashboard/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewLoginUsecase(userRepo domain.UserRepository, jwtSecret string) *loginUsecase {
	return &loginUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *loginUsecase) Login(username, password string) (string, *domain.User, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		// Return the error directly (it could be ErrUserNotFound or DB error)
		return "", nil, err
	}

	// Check if user is active
	if user.Status != "active" {
		return "", nil, domain.ErrInvalidPassword // Or create a new error for inactive users
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, domain.ErrInvalidPassword
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	u.userRepo.Update(user)

	// Generate JWT with role information
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role_id":  user.RoleID,
		"role":     user.Role.Name,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
