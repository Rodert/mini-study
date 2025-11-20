package model

import "time"

// TableName 指定表名
func (Banner) TableName() string {
	return "banners"
}

// Banner represents a carousel banner configured by admin.
type Banner struct {
	Base
	Title        string     `gorm:"size:255;comment:标题" json:"title"`
	ImageURL     string     `gorm:"size:512;comment:图片URL" json:"image_url"`
	LinkURL      string     `gorm:"size:512;comment:跳转链接URL" json:"link_url"`
	VisibleRoles string     `gorm:"size:16;default:'both';comment:可见角色(employee员工/manager店长/both全部)" json:"visible_roles"`
	SortOrder    int        `gorm:"default:0;comment:排序顺序" json:"sort_order"`
	Status       bool       `gorm:"default:true;comment:状态(启用/禁用)" json:"status"`
	StartAt      *time.Time `gorm:"comment:开始时间" json:"start_at"`
	EndAt        *time.Time `gorm:"comment:结束时间" json:"end_at"`
}

