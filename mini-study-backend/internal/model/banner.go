package model

import "time"

// Banner represents a carousel banner configured by admin.
type Banner struct {
	Base
	Title        string     `gorm:"size:255" json:"title"`
	ImageURL     string     `gorm:"size:512" json:"image_url"`
	LinkURL      string     `gorm:"size:512" json:"link_url"`
	VisibleRoles string     `gorm:"size:16;default:'both'" json:"visible_roles"`
	SortOrder    int        `gorm:"default:0" json:"sort_order"`
	Status       bool       `gorm:"default:true" json:"status"`
	StartAt      *time.Time `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
}

