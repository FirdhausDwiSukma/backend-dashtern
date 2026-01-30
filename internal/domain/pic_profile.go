package domain

import "time"

// PICProfile represents Person In Charge (mentor) profile data
type PICProfile struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Position  string    `json:"position"`
	Division  string    `json:"division"`
	Expertise string    `json:"expertise"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for PICProfile model
func (PICProfile) TableName() string {
	return "pic_profiles"
}
