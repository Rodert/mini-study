package dto

import "time"

// PointTransactionListQuery controls pagination for point records.
type PointTransactionListQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

// PointTransactionResponse represents a single point change.
type PointTransactionResponse struct {
	ID          uint      `json:"id" example:"1"`
	UserID      uint      `json:"user_id" example:"10"`
	Change      int64     `json:"change" example:"1"`
	Source      string    `json:"source" example:"content_completion"`
	ReferenceID string    `json:"reference_id" example:"content:12"`
	ContentID   *uint     `json:"content_id,omitempty" example:"12"`
	Description string    `json:"description" example:"完成课程《产品介绍》"`
	Memo        string    `json:"memo"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserPointDetailResponse aggregates total points and transactions for a user.
type UserPointDetailResponse struct {
	User         UserResponse               `json:"user"`
	TotalPoints  int64                      `json:"total_points"`
	Transactions []PointTransactionResponse `json:"transactions"`
	Pagination   Pagination                 `json:"pagination"`
}

// AdminListPointsQuery filters admin points list.
type AdminListPointsQuery struct {
	Keyword  string `form:"keyword" binding:"omitempty,max=100" example:"张三"`                         // 关键词（工号/姓名/手机号）
	Role     string `form:"role" binding:"omitempty,oneof=employee manager admin" example:"employee"` // 角色过滤
	Page     int    `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

// UserPointListItem represents a user with their point total for list display.
type UserPointListItem struct {
	UserResponse
	Points int64 `json:"points" example:"100"` // 积分总数
}

// AdminListPointsResponse returns paginated list of users with points.
type AdminListPointsResponse struct {
	Items      []UserPointListItem `json:"items"`
	Pagination Pagination          `json:"pagination"`
}
