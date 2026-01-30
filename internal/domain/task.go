package domain

import "time"

// Task represents work assignments for interns
type Task struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	InternID     uint       `gorm:"not null" json:"intern_id"`
	Intern       User       `gorm:"foreignKey:InternID" json:"intern"`
	Title        string     `gorm:"not null" json:"title"`
	Description  string     `json:"description"`
	Status       string     `gorm:"default:todo" json:"status"` // todo, in_progress, done
	QualityScore *int       `json:"quality_score"`              // 0-100
	Deadline     time.Time  `json:"deadline"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

// TableName specifies the table name for Task model
func (Task) TableName() string {
	return "tasks"
}
