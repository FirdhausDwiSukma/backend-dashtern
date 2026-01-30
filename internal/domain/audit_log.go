package domain

import "time"

// AuditLog represents system audit trail for tracking user actions
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Action     string    `gorm:"not null" json:"action"` // created, updated, deleted, login
	EntityType string    `json:"entity_type"`            // user, task, attendance, etc.
	EntityID   *uint     `json:"entity_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}
