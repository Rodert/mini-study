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

// GrowthHandler 处理成长圈相关接口。
type GrowthHandler struct {
	service *service.GrowthService
}

// NewGrowthHandler 创建成长圈处理器。
func NewGrowthHandler(svc *service.GrowthService) *GrowthHandler {
	return &GrowthHandler{service: svc}
}

// ListPublicPosts godoc
// @Summary 查询成长圈动态
// @Description 返回所有已审核通过的成长圈动态列表，可按关键词搜索
// @Tags 成长圈
// @Security Bearer
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} utils.Response{data=[]dto.GrowthPostResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/growth [get]
func (h *GrowthHandler) ListPublicPosts(c *gin.Context) {
	if middleware.GetUserID(c) == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var query dto.GrowthListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	posts, err := h.service.ListPublic(query.Keyword)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(posts).JSON(c)
}

// ListMyPosts godoc
// @Summary 查询我的成长圈动态
// @Description 返回当前登录用户发布的成长圈动态列表，可按状态和关键词筛选
// @Tags 成长圈
// @Security Bearer
// @Produce json
// @Param status query string false "状态过滤 pending/approved/rejected"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} utils.Response{data=[]dto.GrowthPostResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/growth/mine [get]
func (h *GrowthHandler) ListMyPosts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var query dto.GrowthMyListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	posts, err := h.service.ListMine(userID, query)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(posts).JSON(c)
}

// CreatePost godoc
// @Summary 店长发布成长圈动态
// @Description 仅店长可以发布成长圈动态，支持文本+多图
// @Tags 成长圈
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.CreateGrowthPostRequest true "成长圈动态"
// @Success 200 {object} utils.Response{data=dto.GrowthPostResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/growth [post]
func (h *GrowthHandler) CreatePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var req dto.CreateGrowthPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	post, err := h.service.CreatePost(userID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(post).JSON(c)
}

// DeletePost godoc
// @Summary 删除成长圈动态
// @Description 店长可删除自己未通过审核的动态，管理员可删除任意动态
// @Tags 成长圈
// @Security Bearer
// @Produce json
// @Param id path int true "动态ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/growth/{id} [delete]
func (h *GrowthHandler) DeletePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || postID == 0 {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的动态ID").JSON(c)
		return
	}

	if err := h.service.Delete(userID, uint(postID)); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(nil).JSON(c)
}

// AdminListPosts godoc
// @Summary 管理员查询成长圈动态列表
// @Description 管理员可按状态和关键词筛选成长圈动态
// @Tags 管理后台-成长圈
// @Security Bearer
// @Produce json
// @Param status query string false "状态过滤 pending/approved/rejected"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} utils.Response{data=[]dto.GrowthPostResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/growth [get]
func (h *GrowthHandler) AdminListPosts(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var query dto.AdminGrowthListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	posts, err := h.service.AdminList(adminID, query)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(posts).JSON(c)
}

// AdminApprovePost godoc
// @Summary 管理员审核通过成长圈动态
// @Description 将指定成长圈动态状态设置为 approved
// @Tags 管理后台-成长圈
// @Security Bearer
// @Produce json
// @Param id path int true "动态ID"
// @Success 200 {object} utils.Response{data=dto.GrowthPostResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/growth/{id}/approve [post]
func (h *GrowthHandler) AdminApprovePost(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || postID == 0 {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的动态ID").JSON(c)
		return
	}

	post, err := h.service.Approve(adminID, uint(postID))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(post).JSON(c)
}

// AdminRejectPost godoc
// @Summary 管理员拒绝成长圈动态
// @Description 将指定成长圈动态状态设置为 rejected
// @Tags 管理后台-成长圈
// @Security Bearer
// @Produce json
// @Param id path int true "动态ID"
// @Success 200 {object} utils.Response{data=dto.GrowthPostResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/growth/{id}/reject [post]
func (h *GrowthHandler) AdminRejectPost(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || postID == 0 {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的动态ID").JSON(c)
		return
	}

	post, err := h.service.Reject(adminID, uint(postID))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(post).JSON(c)
}
