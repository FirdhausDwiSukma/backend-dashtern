package domain

import "time"

// InternProfile represents internship-specific profile data
type InternProfile struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	PICID     uint      `gorm:"not null" json:"pic_id"`
	PIC       User      `gorm:"foreignKey:PICID" json:"pic"`
	Batch     string    `json:"batch"`
	Division  string    `json:"division"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Education string    `json:"education"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for InternProfile model
func (InternProfile) TableName() string {
	return "intern_profiles"
}
