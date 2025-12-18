package dto

import "time"

// CreateGrowthPostRequest 店长创建成长圈动态请求体。
type CreateGrowthPostRequest struct {
	Content    string   `json:"content" binding:"required,min=1,max=1000" example:"今天门店销售突破50万，大家辛苦了！"` // 动态文本内容
	ImagePaths []string `json:"image_paths" binding:"omitempty,max=9,dive,required"`                                 // 图片路径数组，最多9张
}

// GrowthListQuery 成长圈公开列表查询参数。
type GrowthListQuery struct {
	Keyword string `form:"keyword" example:"销售"` // 搜索关键词，按内容模糊匹配
}

// GrowthMyListQuery 当前用户自己的成长圈列表查询参数。
type GrowthMyListQuery struct {
	Keyword string `form:"keyword" example:"培训"`                                                        // 搜索关键词
	Status  string `form:"status" binding:"omitempty,oneof=pending approved rejected" example:"pending"` // 状态过滤
}

// AdminGrowthListQuery 管理员成长圈列表查询参数。
type AdminGrowthListQuery struct {
	Keyword string `form:"keyword" example:"奖励"`                                                        // 搜索关键词
	Status  string `form:"status" binding:"omitempty,oneof=pending approved rejected" example:"pending"` // 状态过滤
}

// GrowthPostResponse 成长圈动态返回结构。
type GrowthPostResponse struct {
	ID            uint       `json:"id" example:"1"`                         // 动态ID
	Content       string     `json:"content" example:"今天门店销售突破50万，大家辛苦了！"` // 文本内容
	ImagePaths    []string   `json:"image_paths,omitempty"`                   // 图片路径数组（相对路径），前端使用 buildFileUrl 转全路径
	Status        string     `json:"status" example:"approved"`            // 状态：pending/approved/rejected
	PublisherID   uint       `json:"publisher_id" example:"3"`             // 发布者用户ID
	PublisherName string     `json:"publisher_name" example:"张店长"`       // 发布者姓名
	PublisherRole string     `json:"publisher_role" example:"manager"`     // 发布者角色：employee/manager/admin
	CreatedAt     time.Time  `json:"created_at" example:"2024-01-01T12:00:00Z"` // 创建时间
	ApprovedAt    *time.Time `json:"approved_at,omitempty" example:"2024-01-01T13:00:00Z"`   // 审核通过时间
}
