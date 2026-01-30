package domain

import "time"

// Attendance represents daily attendance records for interns
type Attendance struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	InternID    uint       `gorm:"not null" json:"intern_id"`
	Intern      User       `gorm:"foreignKey:InternID" json:"intern"`
	Date        time.Time  `gorm:"not null" json:"date"`
	Status      string     `gorm:"not null" json:"status"` // hadir, izin, alpha
	CheckInTime *time.Time `json:"check_in_time"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TableName specifies the table name for Attendance model
func (Attendance) TableName() string {
	return "attendance"
}
