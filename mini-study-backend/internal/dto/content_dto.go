package dto

import "time"

// AdminCreateContentRequest defines payload to create content.
type AdminCreateContentRequest struct {
	Title           string `json:"title" binding:"required,min=1,max=255" example:"产品培训视频"`                      // 内容标题
	Type            string `json:"type" binding:"required,oneof=doc video" example:"video"`                      // 内容类型：doc(文档) 或 video(视频)
	CategoryID      uint   `json:"category_id" binding:"required" example:"1"`                                   // 分类ID
	FilePath        string `json:"file_path" binding:"required" example:"/uploads/video.mp4"`                    // 文件存储路径
	CoverURL        string `json:"cover_url" example:"https://example.com/cover.jpg"`                            // 封面图片URL
	Summary         string `json:"summary" example:"本视频介绍产品核心功能"`                                                // 内容摘要
	VisibleRoles    string `json:"visible_roles" binding:"omitempty,oneof=employee manager both" example:"both"` // 可见角色：employee(员工) manager(店长) both(全部)
	Status          string `json:"status" binding:"omitempty,oneof=draft published" example:"published"`         // 状态：draft(草稿) published(已发布)
	DurationSeconds int64  `json:"duration_seconds" example:"3600"`                                              // 视频时长（秒）
}

// AdminUpdateContentRequest defines payload to update content.
type AdminUpdateContentRequest struct {
	Title           string `json:"title" binding:"omitempty,min=1,max=255" example:"产品培训视频(更新)"`                 // 内容标题
	Type            string `json:"type" binding:"omitempty,oneof=doc video" example:"video"`                     // 内容类型：doc(文档) 或 video(视频)
	CategoryID      uint   `json:"category_id" example:"1"`                                                      // 分类ID
	FilePath        string `json:"file_path" example:"/uploads/video.mp4"`                                       // 文件存储路径
	CoverURL        string `json:"cover_url" example:"https://example.com/cover.jpg"`                            // 封面图片URL
	Summary         string `json:"summary" example:"本视频介绍产品核心功能"`                                                // 内容摘要
	VisibleRoles    string `json:"visible_roles" binding:"omitempty,oneof=employee manager both" example:"both"` // 可见角色：employee(员工) manager(店长) both(全部)
	Status          string `json:"status" binding:"omitempty,oneof=draft published offline" example:"published"` // 状态：draft(草稿) published(已发布) offline(下线)
	DurationSeconds int64  `json:"duration_seconds" example:"3600"`                                              // 视频时长（秒）
}

// AdminListContentRequest filters admin content list.
type AdminListContentRequest struct {
	CategoryID uint   `form:"category_id" example:"1"`                                                      // 分类ID
	Type       string `form:"type" binding:"omitempty,oneof=doc video" example:"video"`                     // 内容类型：doc(文档) 或 video(视频)
	Status     string `form:"status" binding:"omitempty,oneof=draft published offline" example:"published"` // 状态：draft(草稿) published(已发布) offline(下线)
}

// PublishedContentQuery filters public content list.
type PublishedContentQuery struct {
	CategoryID uint   `form:"category_id" example:"1"`                                  // 分类ID
	Type       string `form:"type" binding:"omitempty,oneof=doc video" example:"video"` // 内容类型：doc(文档) 或 video(视频)
}

// ContentCategoryResponse represents category info.
type ContentCategoryResponse struct {
	ID        uint   `json:"id" example:"1"`            // 分类ID
	Name      string `json:"name" example:"产品培训"`       // 分类名称
	RoleScope string `json:"role_scope" example:"both"` // 可见角色范围：employee(员工) manager(店长) both(全部)
	SortOrder int    `json:"sort_order" example:"1"`    // 排序序号
}

// ContentResponse is returned to clients.
type ContentResponse struct {
	ID              uint       `json:"id" example:"1"`                                      // 内容ID
	Title           string     `json:"title" example:"产品培训视频"`                              // 内容标题
	Type            string     `json:"type" example:"video"`                                // 内容类型：doc(文档) 或 video(视频)
	CategoryID      uint       `json:"category_id" example:"1"`                             // 分类ID
	CategoryName    string     `json:"category_name" example:"产品培训"`                        // 分类名称
	FilePath        string     `json:"file_path" example:"/uploads/video.mp4"`              // 文件存储路径
	CoverURL        string     `json:"cover_url" example:"https://example.com/cover.jpg"`   // 封面图片URL
	Summary         string     `json:"summary" example:"本视频介绍产品核心功能"`                       // 内容摘要
	Status          string     `json:"status" example:"published"`                          // 状态：draft(草稿) published(已发布) offline(下线)
	VisibleRoles    string     `json:"visible_roles" example:"both"`                        // 可见角色：employee(员工) manager(店长) both(全部)
	DurationSeconds int64      `json:"duration_seconds" example:"3600"`                     // 视频时长（秒）
	PublishAt       *time.Time `json:"publish_at,omitempty" example:"2024-01-01T00:00:00Z"` // 发布时间
}

// LearningProgressRequest upserts learning progress.
type LearningProgressRequest struct {
	ContentID     uint  `json:"content_id" binding:"required" example:"1"`             // 内容ID
	VideoPosition int64 `json:"video_position" binding:"required,gte=0" example:"120"` // 视频播放位置（秒）
}

// LearningProgressResponse returns progress info.
type LearningProgressResponse struct {
	ContentID       uint   `json:"content_id" example:"1"`          // 内容ID
	VideoPosition   int64  `json:"video_position" example:"120"`    // 视频播放位置（秒）
	DurationSeconds int64  `json:"duration_seconds" example:"3600"` // 视频总时长（秒）
	Progress        int    `json:"progress" example:"3"`            // 学习进度百分比（0-100）
	Status          string `json:"status" example:"in_progress"`    // 学习状态：not_started(未开始) in_progress(进行中) completed(已完成)
}
