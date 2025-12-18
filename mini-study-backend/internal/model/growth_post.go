package model

import "time"

// TableName 指定成长圈表名
func (GrowthPost) TableName() string {
	return "growth_posts"
}

// GrowthPost 成长圈动态，支持文本+多图并记录发布与审核信息。
type GrowthPost struct {
	Base
	CreatorID  uint  `gorm:"not null;index;comment:发布者用户ID" json:"creator_id"`
	Creator    User  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Content    string `gorm:"type:text;not null;comment:动态文本内容" json:"content"`
	ImagePaths string `gorm:"type:text;comment:图片路径数组(JSON)" json:"-"`
	Status     string `gorm:"size:16;default:'pending';comment:状态(pending待审核/approved已通过/rejected已拒绝)" json:"status"`
	ApprovedAt *time.Time `gorm:"comment:审核通过时间" json:"approved_at,omitempty"`
}
