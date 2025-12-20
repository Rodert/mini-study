package dto

import "time"

type AdminCreateNoticeRequest struct {
	Title    string     `json:"title" binding:"required,min=1,max=255"`
	Content  string     `json:"content" binding:"omitempty,max=2000"`
	ImageURL string     `json:"image_url" binding:"omitempty"`
	Status   *bool      `json:"status"`
	StartAt  *time.Time `json:"start_at"`
	EndAt    *time.Time `json:"end_at"`
}

type AdminUpdateNoticeRequest struct {
	Title    string     `json:"title" binding:"omitempty,min=1,max=255"`
	Content  string     `json:"content" binding:"omitempty,max=2000"`
	ImageURL string     `json:"image_url" binding:"omitempty"`
	Status   *bool      `json:"status"`
	StartAt  *time.Time `json:"start_at"`
	EndAt    *time.Time `json:"end_at"`
}

type AdminListNoticeQuery struct {
	Status *bool `form:"status"`
}

type NoticeResponse struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	ImageURL  string     `json:"image_url"`
	Status    bool       `json:"status"`
	StartAt   *time.Time `json:"start_at"`
	EndAt     *time.Time `json:"end_at"`
	CreatedAt time.Time  `json:"created_at"`
}
