package model

import "time"

// TableName 指定表名
func (ContentCategory) TableName() string {
	return "content_categories"
}

// ContentCategory 学习内容分类，固定区分员工/店长可见范围。
type ContentCategory struct {
	Base
	Name      string `gorm:"size:64;uniqueIndex;comment:分类名称" json:"name"`
	RoleScope string `gorm:"size:16;comment:可见角色范围(employee员工/manager店长)" json:"role_scope"`
	SortOrder int    `gorm:"default:0;comment:排序顺序" json:"sort_order"`
	Status    bool   `gorm:"default:true;comment:状态(启用/禁用)" json:"status"`
	CoverURL  string `gorm:"size:512;comment:分类封面图片URL" json:"cover_url"`
}

// TableName 指定表名
func (Content) TableName() string {
	return "contents"
}

// Content 学习内容，支持文档/视频并记录发布信息。
type Content struct {
	Base
	Title           string          `gorm:"size:255;comment:标题" json:"title"`
	Type            string          `gorm:"size:16;comment:内容类型(document文档/video视频/article图文)" json:"type"`
	CategoryID      uint            `gorm:"comment:分类ID" json:"category_id"`
	Category        ContentCategory `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	VisibleRoles    string          `gorm:"size:16;comment:可见角色(employee员工/manager店长/both全部)" json:"visible_roles"`
	FilePath        string          `gorm:"size:512;comment:文件路径" json:"file_path"`
	CoverURL        string          `gorm:"size:512;comment:封面图片URL" json:"cover_url"`
	Summary         string          `gorm:"type:text;comment:摘要" json:"summary"`
	BodyBlocksJSON  string          `gorm:"type:longtext;comment:图文内容结构(JSON)" json:"-"`
	Status          string          `gorm:"size:16;comment:状态(draft草稿/published已发布)" json:"status"`
	PublishAt       *time.Time      `gorm:"comment:发布时间" json:"publish_at"`
	CreatorID       uint            `gorm:"comment:创建者ID" json:"creator_id"`
	DurationSeconds int64           `gorm:"comment:时长(秒)" json:"duration_seconds"`
}

// TableName 指定表名
func (LearningRecord) TableName() string {
	return "learning_records"
}

// LearningRecord 记录用户在某个内容的学习进度与观看位置。
type LearningRecord struct {
	Base
	UserID        uint       `gorm:"uniqueIndex:idx_user_content;comment:用户ID" json:"user_id"`
	ContentID     uint       `gorm:"uniqueIndex:idx_user_content;comment:内容ID" json:"content_id"`
	Progress      int        `gorm:"default:0;comment:学习进度(百分比)" json:"progress"`
	VideoPosition int64      `gorm:"default:0;comment:视频观看位置(秒)" json:"video_position"`
	Status        string     `gorm:"size:16;comment:状态(learning学习中/completed已完成)" json:"status"`
	CompletedAt   *time.Time `gorm:"comment:完成时间" json:"completed_at"`
}
