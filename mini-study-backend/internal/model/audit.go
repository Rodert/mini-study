package model

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLog stores audit info for critical operations.
type AuditLog struct {
	Base
	ActorID uint   `gorm:"comment:操作者ID" json:"actor_id"`
	Action  string `gorm:"size:64;comment:操作动作" json:"action"`
	Target  string `gorm:"size:128;comment:操作目标" json:"target"`
	Payload string `gorm:"type:text;comment:操作载荷(JSON格式)" json:"payload"`
	Result  string `gorm:"size:32;comment:操作结果(success成功/failed失败)" json:"result"`
}
