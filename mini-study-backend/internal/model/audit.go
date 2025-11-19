package model

// AuditLog stores audit info for critical operations.
type AuditLog struct {
	Base
	ActorID uint   `json:"actor_id"`
	Action  string `gorm:"size:64" json:"action"`
	Target  string `gorm:"size:128" json:"target"`
	Payload string `gorm:"type:text" json:"payload"`
	Result  string `gorm:"size:32" json:"result"`
}
