package dto

// RegisterRequest represents the registration payload.
// manager_ids 中放的是店长的 work_no，而不是数值 ID。
type RegisterRequest struct {
	WorkNo         string   `json:"work_no" binding:"required,min=2,max=50"`
	Phone          string   `json:"phone" binding:"omitempty,max=20"`
	Name           string   `json:"name" binding:"omitempty,max=100"`
	Password       string   `json:"password" binding:"required,min=6"`
	ManagerWorkNos []string `json:"manager_ids" binding:"omitempty,dive,required"`
}

// LoginRequest represents the login payload.
type LoginRequest struct {
	WorkNo   string `json:"work_no" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is returned to the client.
type UserResponse struct {
	ID     uint   `json:"id"`
	WorkNo string `json:"work_no"`
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Status bool   `json:"status"`
}

// TokenResponse contains JWT pair.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequest carries the refresh token for renewal.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateProfileRequest is used by user to update own basic info.
type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"omitempty,max=100"`
	Phone string `json:"phone" binding:"omitempty,max=20"`
}

// AdminCreateManagerRequest is used by admin to create a new manager.
type AdminCreateManagerRequest struct {
	WorkNo   string `json:"work_no" binding:"required,min=2,max=50"`
	Name     string `json:"name" binding:"required,max=100"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

// AdminUpdateEmployeeManagersRequest updates manager bindings for an employee.
type AdminUpdateEmployeeManagersRequest struct {
	ManagerWorkNos []string `json:"manager_ids" binding:"omitempty,dive,required"`
}
