package dto

import "time"

type ArticleBlock struct {
	Type      string `json:"type"`
	Text      string `json:"text,omitempty"`
	ImagePath string `json:"image_path,omitempty"`
}

// AdminCreateContentRequest defines payload to create content.
type AdminCreateContentRequest struct {
	Title           string         `json:"title" binding:"required,min=1,max=255" example:"产品培训视频"`                      // 内容标题
	Type            string         `json:"type" binding:"required,oneof=doc video article" example:"video"`                 // 内容类型：doc(文档) 或 video(视频) 或 article(图文)
	CategoryID      uint           `json:"category_id" binding:"required" example:"1"`                                     // 分类ID
	FilePath        string         `json:"file_path" example:"/uploads/video.mp4"`                                          // 文件存储路径
	CoverURL        string         `json:"cover_url" example:"https://example.com/cover.jpg"`                               // 封面图片URL
	Summary         string         `json:"summary" example:"本视频介绍产品核心功能"`                                                // 内容摘要
	VisibleRoles    string         `json:"visible_roles" binding:"omitempty,oneof=employee manager both" example:"both"`   // 可见角色：employee(员工) manager(店长) both(全部)
	Status          string         `json:"status" binding:"omitempty,oneof=draft published" example:"published"`           // 状态：draft(草稿) published(已发布)
	DurationSeconds int64          `json:"duration_seconds" example:"3600"`                                                 // 视频时长（秒）
	ArticleBlocks   []ArticleBlock `json:"article_blocks,omitempty"`                                                          // 图文内容块，仅在 type=article 时使用
}

// AdminUpdateContentRequest defines payload to update content.
type AdminUpdateContentRequest struct {
	Title           string         `json:"title" binding:"omitempty,min=1,max=255" example:"产品培训视频(更新)"`           // 内容标题
	Type            string         `json:"type" binding:"omitempty,oneof=doc video article" example:"video"`             // 内容类型：doc(文档) 或 video(视频) 或 article(图文)
	CategoryID      uint           `json:"category_id" example:"1"`                                                       // 分类ID
	FilePath        string         `json:"file_path" example:"/uploads/video.mp4"`                                        // 文件存储路径
	CoverURL        string         `json:"cover_url" example:"https://example.com/cover.jpg"`                             // 封面图片URL
	Summary         string         `json:"summary" example:"本视频介绍产品核心功能"`                                                // 内容摘要
	VisibleRoles    string         `json:"visible_roles" binding:"omitempty,oneof=employee manager both" example:"both"` // 可见角色：employee(员工) manager(店长) both(全部)
	Status          string         `json:"status" binding:"omitempty,oneof=draft published offline" example:"published"` // 状态：draft(草稿) published(已发布) offline(下线)
	DurationSeconds int64          `json:"duration_seconds" example:"3600"`                                                // 视频时长（秒）
	ArticleBlocks   []ArticleBlock `json:"article_blocks,omitempty"`                                                        // 图文内容块，仅在 type=article 时使用
}

// AdminListContentRequest filters admin content list.
type AdminListContentRequest struct {
	CategoryID uint   `form:"category_id" example:"1"`                                                      // 分类ID
	Type       string `form:"type" binding:"omitempty,oneof=doc video article" example:"video"`            // 内容类型：doc(文档) 或 video(视频) 或 article(图文)
	Status     string `form:"status" binding:"omitempty,oneof=draft published offline" example:"published"` // 状态：draft(草稿) published(已发布) offline(下线)
}

// PublishedContentQuery filters public content list.
type PublishedContentQuery struct {
	CategoryID uint   `form:"category_id" example:"1"`                                         // 分类ID
	Type       string `form:"type" binding:"omitempty,oneof=doc video article" example:"video"` // 内容类型：doc(文档) 或 video(视频) 或 article(图文)
}

// ContentCategoryResponse represents category info.
type ContentCategoryResponse struct {
	ID        uint   `json:"id" example:"1"`            // 分类ID
	Name      string `json:"name" example:"产品培训"`       // 分类名称
	RoleScope string `json:"role_scope" example:"both"` // 可见角色范围：employee(员工) manager(店长) both(全部)
	SortOrder int    `json:"sort_order" example:"1"`    // 排序序号
	CoverURL  string `json:"cover_url" example:"/uploads/category-cover.png"` // 分类封面图片URL
	Count     int64  `json:"count" example:"5"`         // 该分类下已发布且对当前用户可见的课程数量
}

type AdminUpdateCategoryRequest struct {
	CoverURL *string `json:"cover_url"`
}

// ContentResponse is returned to clients.
type ContentResponse struct {
	ID              uint           `json:"id" example:"1"`                                      // 内容ID
	Title           string         `json:"title" example:"产品培训视频"`                              // 内容标题
	Type            string         `json:"type" example:"video"`                                // 内容类型：doc(文档) 或 video(视频) 或 article(图文)
	CategoryID      uint           `json:"category_id" example:"1"`                             // 分类ID
	CategoryName    string         `json:"category_name" example:"产品培训"`                        // 分类名称
	FilePath        string         `json:"file_path" example:"/uploads/video.mp4"`              // 文件存储路径
	CoverURL        string         `json:"cover_url" example:"https://example.com/cover.jpg"`   // 封面图片URL
	Summary         string         `json:"summary" example:"本视频介绍产品核心功能"`                       // 内容摘要
	Status          string         `json:"status" example:"published"`                          // 状态：draft(草稿) published(已发布) offline(下线)
	VisibleRoles    string         `json:"visible_roles" example:"both"`                        // 可见角色：employee(员工) manager(店长) both(全部)
	DurationSeconds int64          `json:"duration_seconds" example:"3600"`                     // 视频时长（秒）
	PublishAt       *time.Time     `json:"publish_at,omitempty" example:"2024-01-01T00:00:00Z"` // 发布时间
	ArticleBlocks   []ArticleBlock `json:"article_blocks,omitempty"`                             // 图文内容块，仅在 type=article 时返回
}

// LearningProgressRequest upserts learning progress.
type LearningProgressRequest struct {
	ContentID     uint  `json:"content_id" binding:"required" example:"1"`              // 内容ID
	VideoPosition int64 `json:"video_position" binding:"gte=0" example:"120"`           // 视频播放位置（秒），0表示文档类型或初始状态
}

// LearningProgressResponse returns progress info.
type LearningProgressResponse struct {
	ContentID       uint   `json:"content_id" example:"1"`          // 内容ID
	VideoPosition   int64  `json:"video_position" example:"120"`    // 视频播放位置（秒）
	DurationSeconds int64  `json:"duration_seconds" example:"3600"` // 视频总时长（秒）
	Progress        int    `json:"progress" example:"3"`            // 学习进度百分比（0-100）
	Status          string `json:"status" example:"in_progress"`    // 学习状态：not_started(未开始) in_progress(进行中) completed(已完成)
}

// UserLearningStatsResponse returns learning statistics for a user.
type UserLearningStatsResponse struct {
	UserID         uint    `json:"user_id" example:"1"`                    // 用户ID
	CompletedCount int64   `json:"completed_count" example:"5"`            // 已完成的学习内容数量
	TotalCount     int64   `json:"total_count" example:"10"`              // 已开始学习的内容总数
	TotalContents  int64   `json:"total_contents" example:"20"`           // 该角色可见的已发布内容总数
	CompletionRate float64 `json:"completion_rate" example:"25.0"`        // 完成率（百分比）
}

// ContentCompletionStatsResponse returns completion statistics for a content.
type ContentCompletionStatsResponse struct {
	ContentID      uint    `json:"content_id" example:"1"`                // 内容ID
	ContentTitle   string  `json:"content_title" example:"产品培训视频"`      // 内容标题
	CompletedCount int64   `json:"completed_count" example:"15"`           // 已完成学习的用户数量
	TotalCount     int64   `json:"total_count" example:"30"`              // 已开始学习的用户总数
	CompletionRate float64 `json:"completion_rate" example:"50.0"`        // 完成率（百分比）
}
