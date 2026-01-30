package domain

import "time"

// MentorReview represents mentor evaluations of intern potential
type MentorReview struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	InternID        uint      `gorm:"not null" json:"intern_id"`
	Intern          User      `gorm:"foreignKey:InternID" json:"intern"`
	PICID           uint      `gorm:"not null" json:"pic_id"`
	PIC             User      `gorm:"foreignKey:PICID" json:"pic"`
	LearningAbility int       `gorm:"not null" json:"learning_ability"` // 1-5
	Initiative      int       `gorm:"not null" json:"initiative"`       // 1-5
	Communication   int       `gorm:"not null" json:"communication"`    // 1-5
	ProblemSolving  int       `gorm:"not null" json:"problem_solving"`  // 1-5
	Notes           string    `json:"notes"`
	Period          string    `gorm:"not null" json:"period"` // Format: 2026-01
	CreatedAt       time.Time `json:"created_at"`
}

// TableName specifies the table name for MentorReview model
func (MentorReview) TableName() string {
	return "mentor_reviews"
}
