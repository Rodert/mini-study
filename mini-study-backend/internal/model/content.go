package model

import "time"

// ContentCategory 学习内容分类，固定区分员工/店长可见范围。
type ContentCategory struct {
	Base
	Name      string `gorm:"size:64;uniqueIndex" json:"name"`
	RoleScope string `gorm:"size:16" json:"role_scope"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	Status    bool   `gorm:"default:true" json:"status"`
}

// Content 学习内容，支持文档/视频并记录发布信息。
type Content struct {
	Base
	Title           string          `gorm:"size:255" json:"title"`
	Type            string          `gorm:"size:16" json:"type"`
	CategoryID      uint            `json:"category_id"`
	Category        ContentCategory `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	VisibleRoles    string          `gorm:"size:16" json:"visible_roles"`
	FilePath        string          `gorm:"size:512" json:"file_path"`
	CoverURL        string          `gorm:"size:512" json:"cover_url"`
	Summary         string          `gorm:"type:text" json:"summary"`
	Status          string          `gorm:"size:16" json:"status"`
	PublishAt       *time.Time      `json:"publish_at"`
	CreatorID       uint            `json:"creator_id"`
	DurationSeconds int64           `json:"duration_seconds"`
}

// LearningRecord 记录用户在某个内容的学习进度与观看位置。
type LearningRecord struct {
	Base
	UserID        uint       `gorm:"uniqueIndex:idx_user_content" json:"user_id"`
	ContentID     uint       `gorm:"uniqueIndex:idx_user_content" json:"content_id"`
	Progress      int        `gorm:"default:0" json:"progress"`
	VideoPosition int64      `gorm:"default:0" json:"video_position"`
	Status        string     `gorm:"size:16" json:"status"`
	CompletedAt   *time.Time `json:"completed_at"`
}
