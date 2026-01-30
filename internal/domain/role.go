package domain

import "time"

// Role represents user role types in the system
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"` // super_admin, hr, pic, intern
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "roles"
}
