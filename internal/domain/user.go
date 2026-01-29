package domain

import "gorm.io/gorm"

// User represents the user entity
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Don't return password in JSON
	Role     string `gorm:"default:user" json:"role"`
	gorm.Model
}

// UserRepository defines the methods that any storage layer must implement
type UserRepository interface {
	GetByUsername(username string) (*User, error)
	Create(user *User) error
}

// UserUsecase defines the methods that the business logic layer must implement
type UserUsecase interface {
	Login(username, password string) (string, error)
}
