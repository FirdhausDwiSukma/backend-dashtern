package domain

import "time"

// PerformanceScore represents calculated performance metrics from tasks and attendance
type PerformanceScore struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	InternID        uint      `gorm:"not null" json:"intern_id"`
	Intern          User      `gorm:"foreignKey:InternID" json:"intern"`
	Period          string    `gorm:"not null" json:"period"` // Format: 2026-01
	TaskScore       float64   `json:"task_score"`
	AttendanceScore float64   `json:"attendance_score"`
	QualityScore    float64   `json:"quality_score"`
	FinalScore      float64   `json:"final_score"` // 0-100
	CreatedAt       time.Time `json:"created_at"`
}

// TableName specifies the table name for PerformanceScore model
func (PerformanceScore) TableName() string {
	return "performance_scores"
}
