package model

// User represents an application user.
type User struct {
	Base
	WorkNo       string `gorm:"size:50;uniqueIndex;not null" json:"work_no"`
	Phone        string `gorm:"size:20" json:"phone"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Role         Role   `gorm:"size:16;default:'employee'" json:"role"`
	Name         string `gorm:"size:100" json:"name"`
	Status       bool   `gorm:"default:true" json:"status"`
}
