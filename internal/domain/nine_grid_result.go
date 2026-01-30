package domain

import "time"

// NineGridResult represents the final 9-grid positioning combining performance and potential
type NineGridResult struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	InternID         uint      `gorm:"not null" json:"intern_id"`
	Intern           User      `gorm:"foreignKey:InternID" json:"intern"`
	Period           string    `gorm:"not null" json:"period"` // Format: 2026-01
	PerformanceScore float64   `json:"performance_score"`
	PotentialScore   float64   `json:"potential_score"`
	PerformanceLevel string    `json:"performance_level"` // low, medium, high
	PotentialLevel   string    `json:"potential_level"`   // low, medium, high
	GridPosition     string    `json:"grid_position"`     // e.g., "high-high", "medium-low"
	Recommendation   string    `json:"recommendation"`
	GeneratedAt      time.Time `json:"generated_at"`
}

// TableName specifies the table name for NineGridResult model
func (NineGridResult) TableName() string {
	return "nine_grid_results"
}
