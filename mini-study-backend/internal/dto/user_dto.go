package dto

// RegisterRequest represents the registration payload.
// manager_ids 中放的是店长的 work_no，而不是数值 ID。
type RegisterRequest struct {
	WorkNo         string   `json:"work_no" binding:"required,min=2,max=50" example:"E001"`            // 工号
	Phone          string   `json:"phone" binding:"omitempty,max=20" example:"13800138000"`            // 手机号
	Name           string   `json:"name" binding:"omitempty,max=100" example:"张三"`                     // 姓名
	Password       string   `json:"password" binding:"required,min=6" example:"123456"`                // 密码（至少6位）
	ManagerWorkNos []string `json:"manager_ids" binding:"omitempty,dive,required" example:"M001,M002"` // 店长工号列表
}

// LoginRequest represents the login payload.
type LoginRequest struct {
	WorkNo   string `json:"work_no" binding:"required" example:"E001"`    // 工号
	Password string `json:"password" binding:"required" example:"123456"` // 密码
}

// UserResponse is returned to the client.
type UserResponse struct {
	ID     uint   `json:"id" example:"1"`              // 用户ID
	WorkNo string `json:"work_no" example:"E001"`      // 工号
	Phone  string `json:"phone" example:"13800138000"` // 手机号
	Name   string `json:"name" example:"张三"`           // 姓名
	Role   string `json:"role" example:"employee"`     // 角色：employee(员工) manager(店长) admin(管理员)
	Status bool   `json:"status" example:"true"`       // 状态：true(启用) false(禁用)
}

// TokenResponse contains JWT pair.
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // 访问令牌
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 刷新令牌
}

// RefreshTokenRequest carries the refresh token for renewal.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 刷新令牌
}

// UpdateProfileRequest is used by user to update own basic info.
type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"omitempty,max=100" example:"张三"`          // 姓名
	Phone string `json:"phone" binding:"omitempty,max=20" example:"13800138000"` // 手机号
}

// AdminCreateManagerRequest is used by admin to create a new manager.
type AdminCreateManagerRequest struct {
	WorkNo   string `json:"work_no" binding:"required,min=2,max=50" example:"M001"` // 工号
	Name     string `json:"name" binding:"required,max=100" example:"李店长"`          // 姓名
	Phone    string `json:"phone" binding:"omitempty,max=20" example:"13800138000"` // 手机号
	Password string `json:"password" binding:"required,min=6" example:"123456"`     // 密码（至少6位）
}

// AdminCreateEmployeeRequest is used by admin to create a new employee.
// manager_ids 中放的是店长的 work_no，而不是数值 ID。
type AdminCreateEmployeeRequest struct {
	WorkNo         string   `json:"work_no" binding:"required,min=2,max=50" example:"E001"`            // 工号
	Name           string   `json:"name" binding:"required,max=100" example:"张三"`                      // 姓名
	Phone          string   `json:"phone" binding:"omitempty,max=20" example:"13800138000"`            // 手机号
	Password       string   `json:"password" binding:"required,min=6" example:"123456"`                // 密码（至少6位）
	ManagerWorkNos []string `json:"manager_ids" binding:"omitempty,dive,required" example:"M001,M002"` // 店长工号列表
}

// AdminUpdateEmployeeManagersRequest updates manager bindings for an employee.
type AdminUpdateEmployeeManagersRequest struct {
	ManagerWorkNos []string `json:"manager_ids" binding:"omitempty,dive,required" example:"M001,M002"` // 店长工号列表
}
