package domain

import "time"

// PotentialScore represents calculated potential metrics from mentor reviews
type PotentialScore struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	InternID       uint      `gorm:"not null" json:"intern_id"`
	Intern         User      `gorm:"foreignKey:InternID" json:"intern"`
	Period         string    `gorm:"not null" json:"period"` // Format: 2026-01
	MentorAvgScore float64   `json:"mentor_avg_score"`       // 0-100
	CreatedAt      time.Time `json:"created_at"`
}

// TableName specifies the table name for PotentialScore model
func (PotentialScore) TableName() string {
	return "potential_scores"
}
