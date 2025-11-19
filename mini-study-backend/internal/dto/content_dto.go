package dto

import "time"

// AdminCreateContentRequest defines payload to create content.
type AdminCreateContentRequest struct {
	Title           string `json:"title" binding:"required,min=1,max=255"`
	Type            string `json:"type" binding:"required,oneof=doc video"`
	CategoryID      uint   `json:"category_id" binding:"required"`
	FilePath        string `json:"file_path" binding:"required"`
	CoverURL        string `json:"cover_url"`
	Summary         string `json:"summary"`
	VisibleRoles    string `json:"visible_roles" binding:"omitempty,oneof=employee manager both"`
	Status          string `json:"status" binding:"omitempty,oneof=draft published"`
	DurationSeconds int64  `json:"duration_seconds"`
}

// AdminUpdateContentRequest defines payload to update content.
type AdminUpdateContentRequest struct {
	Title           string `json:"title" binding:"omitempty,min=1,max=255"`
	Type            string `json:"type" binding:"omitempty,oneof=doc video"`
	CategoryID      uint   `json:"category_id"`
	FilePath        string `json:"file_path"`
	CoverURL        string `json:"cover_url"`
	Summary         string `json:"summary"`
	VisibleRoles    string `json:"visible_roles" binding:"omitempty,oneof=employee manager both"`
	Status          string `json:"status" binding:"omitempty,oneof=draft published offline"`
	DurationSeconds int64  `json:"duration_seconds"`
}

// AdminListContentRequest filters admin content list.
type AdminListContentRequest struct {
	CategoryID uint   `form:"category_id"`
	Type       string `form:"type" binding:"omitempty,oneof=doc video"`
	Status     string `form:"status" binding:"omitempty,oneof=draft published offline"`
}

// PublishedContentQuery filters public content list.
type PublishedContentQuery struct {
	CategoryID uint   `form:"category_id"`
	Type       string `form:"type" binding:"omitempty,oneof=doc video"`
}

// ContentCategoryResponse represents category info.
type ContentCategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	RoleScope string `json:"role_scope"`
	SortOrder int    `json:"sort_order"`
}

// ContentResponse is returned to clients.
type ContentResponse struct {
	ID              uint       `json:"id"`
	Title           string     `json:"title"`
	Type            string     `json:"type"`
	CategoryID      uint       `json:"category_id"`
	CategoryName    string     `json:"category_name"`
	FilePath        string     `json:"file_path"`
	CoverURL        string     `json:"cover_url"`
	Summary         string     `json:"summary"`
	Status          string     `json:"status"`
	VisibleRoles    string     `json:"visible_roles"`
	DurationSeconds int64      `json:"duration_seconds"`
	PublishAt       *time.Time `json:"publish_at,omitempty"`
}

// LearningProgressRequest upserts learning progress.
type LearningProgressRequest struct {
	ContentID     uint  `json:"content_id" binding:"required"`
	VideoPosition int64 `json:"video_position" binding:"required,gte=0"`
}

// LearningProgressResponse returns progress info.
type LearningProgressResponse struct {
	ContentID       uint   `json:"content_id"`
	VideoPosition   int64  `json:"video_position"`
	DurationSeconds int64  `json:"duration_seconds"`
	Progress        int    `json:"progress"`
	Status          string `json:"status"`
}
