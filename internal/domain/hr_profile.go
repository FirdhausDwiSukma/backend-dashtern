package domain

import "time"

// HRProfile represents HR user profile data
type HRProfile struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName specifies the table name for HRProfile model
func (HRProfile) TableName() string {
	return "hr_profiles"
}
