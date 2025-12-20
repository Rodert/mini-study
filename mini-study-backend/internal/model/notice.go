package model

import "time"

// TableName specifies custom table name for Notice.
func (Notice) TableName() string {
	return "notices"
}

// Notice represents a system announcement configured by admin.
type Notice struct {
	Base
	Title    string     `gorm:"size:255;comment:标题" json:"title"`
	Content  string     `gorm:"type:text;comment:内容" json:"content"`
	ImageURL string     `gorm:"size:512;comment:图片URL" json:"image_url"`
	Status   bool       `gorm:"default:true;comment:状态(启用/禁用)" json:"status"`
	StartAt  *time.Time `gorm:"comment:开始时间" json:"start_at"`
	EndAt    *time.Time `gorm:"comment:结束时间" json:"end_at"`
}
