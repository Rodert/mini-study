package dto

import "time"

// AdminCreateBannerRequest payload for creating banner.
type AdminCreateBannerRequest struct {
	Title        string     `json:"title" binding:"required,min=1,max=255" example:"春季促销活动"`                    // 轮播图标题
	ImageURL     string     `json:"image_url" binding:"required,url" example:"https://example.com/banner.jpg"`  // 图片URL
	LinkURL      string     `json:"link_url" binding:"required,url" example:"https://example.com/promotion"`     // 跳转链接
	VisibleRoles string     `json:"visible_roles" binding:"omitempty,oneof=employee manager both" example:"both"` // 可见角色：employee(员工) manager(店长) both(全部)
	SortOrder    int        `json:"sort_order" example:"1"`                                                      // 排序序号
	Status       *bool      `json:"status" example:"true"`                                                      // 是否启用
	StartAt      *time.Time `json:"start_at" example:"2024-01-01T00:00:00Z"`                                      // 开始时间
	EndAt        *time.Time `json:"end_at" example:"2024-12-31T23:59:59Z"`                                        // 结束时间
}

// AdminUpdateBannerRequest payload for updating banner.
type AdminUpdateBannerRequest struct {
	Title        string     `json:"title" binding:"omitempty,min=1,max=255" example:"春季促销活动(更新)"`           // 轮播图标题
	ImageURL     string     `json:"image_url" binding:"omitempty,url" example:"https://example.com/banner.jpg"` // 图片URL
	LinkURL      string     `json:"link_url" binding:"omitempty,url" example:"https://example.com/promotion"`    // 跳转链接
	VisibleRoles string     `json:"visible_roles" binding:"omitempty,oneof=employee manager both" example:"both"` // 可见角色：employee(员工) manager(店长) both(全部)
	SortOrder    *int       `json:"sort_order" example:"1"`                                                      // 排序序号
	Status       *bool      `json:"status" example:"true"`                                                      // 是否启用
	StartAt      *time.Time `json:"start_at" example:"2024-01-01T00:00:00Z"`                                      // 开始时间
	EndAt        *time.Time `json:"end_at" example:"2024-12-31T23:59:59Z"`                                        // 结束时间
}

// AdminListBannerQuery filters admin banner list.
type AdminListBannerQuery struct {
	Status *bool `form:"status" example:"true"` // 是否启用：true(启用) false(禁用)
}

// BannerResponse is returned to client.
type BannerResponse struct {
	ID           uint       `json:"id" example:"1"`                                                      // 轮播图ID
	Title        string     `json:"title" example:"春季促销活动"`                                         // 轮播图标题
	ImageURL     string     `json:"image_url" example:"https://example.com/banner.jpg"`                 // 图片URL
	LinkURL      string     `json:"link_url" example:"https://example.com/promotion"`                   // 跳转链接
	VisibleRoles string     `json:"visible_roles" example:"both"`                                      // 可见角色：employee(员工) manager(店长) both(全部)
	SortOrder    int        `json:"sort_order" example:"1"`                                             // 排序序号
	Status       bool       `json:"status" example:"true"`                                              // 是否启用
	StartAt      *time.Time `json:"start_at" example:"2024-01-01T00:00:00Z"`                            // 开始时间
	EndAt        *time.Time `json:"end_at" example:"2024-12-31T23:59:59Z"`                              // 结束时间
}
