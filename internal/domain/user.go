package domain

import "time"

// User represents the user entity with role-based structure
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	FullName     string     `gorm:"not null" json:"full_name"`
	Username     string     `gorm:"unique;not null" json:"username"`
	Email        string     `gorm:"unique;not null" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"` // Don't return password in JSON
	RoleID       uint       `gorm:"not null" json:"role_id"`
	Role         Role       `gorm:"foreignKey:RoleID" json:"role"`
	Status       string     `gorm:"default:active" json:"status"` // active, inactive
	AvatarURL    *string    `json:"avatar_url"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// UserRepository defines the methods that any storage layer must implement
type UserRepository interface {
	GetByUsername(username string) (*User, error)
	GetByID(id uint) (*User, error)
	GetAll(page, limit int) ([]User, int64, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

// UserUsecase defines the methods that the business logic layer must implement
type UserUsecase interface {
	Login(username, password string) (string, *User, error)
	GetAll(page, limit int) ([]User, int64, error)
	GetByID(id uint) (*User, error)
	Create(fullName, username, email, password string, roleID uint) (*User, error)
	Update(id uint, fullName, email string, status string) (*User, error)
	Deactivate(id uint) error
	HardDelete(id uint) error
}
