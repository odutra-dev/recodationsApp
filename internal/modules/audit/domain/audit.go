package domain

import "time"

// AuditLog tracks mutations and sensitive actions across the system.
type AuditLog struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	UserID     *int64    `json:"userId,omitempty" gorm:"index"`
	Action     string    `json:"action" gorm:"size:50;not null"`
	Resource   string    `json:"resource" gorm:"size:100;not null"`
	ResourceID *string   `json:"resourceId,omitempty" gorm:"size:255"`
	StatusCode int       `json:"statusCode" gorm:"not null"`
	IPAddress  string    `json:"ipAddress" gorm:"size:50"`
	TraceID    string    `json:"traceId" gorm:"size:100;index"`
	Details    string    `json:"details" gorm:"type:text"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
