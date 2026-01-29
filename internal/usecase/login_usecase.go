package usecase

import (
	"backend-dashboard/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewLoginUsecase(userRepo domain.UserRepository, jwtSecret string) domain.UserUsecase {
	return &loginUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *loginUsecase) Login(username, password string) (string, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials") // Don't reveal if user exists or not for security
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
