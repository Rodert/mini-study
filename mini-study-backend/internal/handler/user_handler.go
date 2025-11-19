package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/middleware"
	"github.com/javapub/mini-study/mini-study-backend/internal/service"
	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

// UserHandler exposes user HTTP endpoints.
type UserHandler struct {
	users  *service.UserService
	tokens *service.TokenService
}

// NewUserHandler builds a handler.
func NewUserHandler(users *service.UserService, tokens *service.TokenService) *UserHandler {
	return &UserHandler{users: users, tokens: tokens}
}

// Register godoc
// @Summary 员工注册
// @Description 员工填写工号、姓名、手机号与密码完成注册，可同时选择多个店长工号
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	user, err := h.users.Register(req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp := dto.UserResponse{
		ID:     user.ID,
		WorkNo: user.WorkNo,
		Phone:  user.Phone,
		Name:   user.Name,
		Role:   user.Role,
		Status: user.Status,
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// Login godoc
// @Summary 用户登录
// @Description 使用工号+密码登录，返回访问令牌与刷新令牌
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "登录信息"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	user, err := h.users.Login(req)
	if err != nil {
		utils.NewErrorResponse(http.StatusUnauthorized, err.Error()).JSON(c)
		return
	}

	tokens, err := h.tokens.GeneratePair(user)
	if err != nil {
		utils.NewErrorResponse(http.StatusInternalServerError, err.Error()).JSON(c)
		return
	}

	utils.NewSuccessResponse(tokens).JSON(c)
}

// ListManagers godoc
// @Summary 查询店长列表
// @Description 返回所有启用状态的店长供员工注册或绑定使用
// @Tags 用户
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.UserResponse}
// @Router /api/v1/users/managers [get]
func (h *UserHandler) ListManagers(c *gin.Context) {
	users, err := h.users.ListManagers()
	if err != nil {
		utils.NewErrorResponse(http.StatusInternalServerError, err.Error()).JSON(c)
		return
	}

	var resp []dto.UserResponse
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID:     u.ID,
			WorkNo: u.WorkNo,
			Phone:  u.Phone,
			Name:   u.Name,
			Role:   u.Role,
			Status: u.Status,
		})
	}

	utils.NewSuccessResponse(resp).JSON(c)
}

// RefreshToken godoc
// @Summary 刷新令牌
// @Description 使用 refresh_token 换取新的访问令牌
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body dto.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/users/token/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	tokens, err := h.tokens.Refresh(req.RefreshToken)
	if err != nil {
		utils.NewErrorResponse(http.StatusUnauthorized, err.Error()).JSON(c)
		return
	}

	utils.NewSuccessResponse(tokens).JSON(c)
}

// UpdateProfile godoc
// @Summary 修改个人信息
// @Description 登录用户可更新姓名与手机号
// @Tags 用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.UpdateProfileRequest true "个人信息"
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/users/me/profile [patch]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	user, err := h.users.UpdateProfile(userID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp := dto.UserResponse{
		ID:     user.ID,
		WorkNo: user.WorkNo,
		Phone:  user.Phone,
		Name:   user.Name,
		Role:   user.Role,
		Status: user.Status,
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// AdminCreateManager godoc
// @Summary 管理员创建店长
// @Description 只有管理员可以创建新的店长账号
// @Tags 管理后台-用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.AdminCreateManagerRequest true "店长信息"
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/managers [post]
func (h *UserHandler) AdminCreateManager(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var req dto.AdminCreateManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	user, err := h.users.CreateManager(adminID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp := dto.UserResponse{
		ID:     user.ID,
		WorkNo: user.WorkNo,
		Phone:  user.Phone,
		Name:   user.Name,
		Role:   user.Role,
		Status: user.Status,
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// AdminPromoteToManager godoc
// @Summary 管理员升级店长
// @Description 将指定员工升级为店长角色
// @Tags 管理后台-用户
// @Security Bearer
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/users/{id}/promote-manager [post]
func (h *UserHandler) AdminPromoteToManager(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || targetID == 0 {
		utils.NewErrorResponse(http.StatusBadRequest, "无效的用户ID").JSON(c)
		return
	}

	user, err := h.users.PromoteToManager(adminID, uint(targetID))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp := dto.UserResponse{
		ID:     user.ID,
		WorkNo: user.WorkNo,
		Phone:  user.Phone,
		Name:   user.Name,
		Role:   user.Role,
		Status: user.Status,
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// AdminUpdateEmployeeManagers godoc
// @Summary 管理员维护员工店长绑定
// @Description 覆盖式更新某个员工绑定的店长工号列表
// @Tags 管理后台-用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Param body body dto.AdminUpdateEmployeeManagersRequest true "店长工号列表"
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/users/{id}/managers [put]
func (h *UserHandler) AdminUpdateEmployeeManagers(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || targetID == 0 {
		utils.NewErrorResponse(http.StatusBadRequest, "无效的用户ID").JSON(c)
		return
	}

	var req dto.AdminUpdateEmployeeManagersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	user, err := h.users.UpdateEmployeeManagers(adminID, uint(targetID), req.ManagerWorkNos)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp := dto.UserResponse{
		ID:     user.ID,
		WorkNo: user.WorkNo,
		Phone:  user.Phone,
		Name:   user.Name,
		Role:   user.Role,
		Status: user.Status,
	}
	utils.NewSuccessResponse(resp).JSON(c)
}
