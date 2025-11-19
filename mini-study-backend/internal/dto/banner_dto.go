package dto

import "time"

// AdminCreateBannerRequest payload for creating banner.
type AdminCreateBannerRequest struct {
	Title        string     `json:"title" binding:"required,min=1,max=255"`
	ImageURL     string     `json:"image_url" binding:"required,url"`
	LinkURL      string     `json:"link_url" binding:"required,url"`
	VisibleRoles string     `json:"visible_roles" binding:"omitempty,oneof=employee manager both"`
	SortOrder    int        `json:"sort_order"`
	Status       *bool      `json:"status"`
	StartAt      *time.Time `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
}

// AdminUpdateBannerRequest payload for updating banner.
type AdminUpdateBannerRequest struct {
	Title        string     `json:"title" binding:"omitempty,min=1,max=255"`
	ImageURL     string     `json:"image_url" binding:"omitempty,url"`
	LinkURL      string     `json:"link_url" binding:"omitempty,url"`
	VisibleRoles string     `json:"visible_roles" binding:"omitempty,oneof=employee manager both"`
	SortOrder    *int       `json:"sort_order"`
	Status       *bool      `json:"status"`
	StartAt      *time.Time `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
}

// AdminListBannerQuery filters admin banner list.
type AdminListBannerQuery struct {
	Status *bool `form:"status"`
}

// BannerResponse is returned to client.
type BannerResponse struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	ImageURL     string     `json:"image_url"`
	LinkURL      string     `json:"link_url"`
	VisibleRoles string     `json:"visible_roles"`
	SortOrder    int        `json:"sort_order"`
	Status       bool       `json:"status"`
	StartAt      *time.Time `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
}
