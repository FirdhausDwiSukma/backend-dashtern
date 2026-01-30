package domain

import "time"

// InternProfile represents internship-specific profile data
type InternProfile struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	PICID      uint      `gorm:"not null" json:"pic_id"`
	PIC        User      `gorm:"foreignKey:PICID" json:"pic"`
	Batch      string    `json:"batch"`
	Division   string    `json:"division"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	University string    `json:"university"`
	Major      string    `json:"major"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName specifies the table name for InternProfile model
func (InternProfile) TableName() string {
	return "intern_profiles"
}

// InternRepository interface
type InternRepository interface {
	Create(userID, picID uint, batch, division, university, major string, startDate, endDate time.Time) (*InternProfile, error)
	GetByID(id uint) (*InternProfile, error)
	GetByUserID(userID uint) (*InternProfile, error)
	GetAll(page, limit int) ([]InternProfile, int64, error)
	Update(id uint, batch, division, university, major string) (*InternProfile, error)
}

// InternUsecase interface
type InternUsecase interface {
	CreateIntern(fullName, username, email, password string, picID uint, batch, division, university, major string, startDate, endDate time.Time) (*User, *InternProfile, error)
	GetInternByID(id uint) (*InternProfile, error)
	GetAllInterns(page, limit int) ([]InternProfile, int64, error)
}
