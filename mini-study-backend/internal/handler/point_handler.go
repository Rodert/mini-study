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

// PointHandler exposes point-related endpoints.
type PointHandler struct {
	points *service.PointService
}

// NewPointHandler builds a PointHandler.
func NewPointHandler(points *service.PointService) *PointHandler {
	return &PointHandler{points: points}
}

// AdminGetUserPoints godoc
// @Summary 管理员查看用户积分明细
// @Description 返回指定用户的积分总数及明细，支持分页
// @Tags 管理后台-积分
// @Security Bearer
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码，从1开始"
// @Param page_size query int false "每页数量，默认20，最大100"
// @Success 200 {object} utils.Response{data=dto.UserPointDetailResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/users/{id}/points [get]
func (h *PointHandler) AdminGetUserPoints(c *gin.Context) {
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

	var query dto.PointTransactionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	detail, err := h.points.AdminUserPointDetails(adminID, uint(targetID), query)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	utils.NewSuccessResponse(detail).JSON(c)
}

// AdminListAllPoints godoc
// @Summary 管理员查看所有用户积分列表
// @Description 返回所有用户的积分列表，按积分降序排列，支持分页和搜索
// @Tags 管理后台-积分
// @Security Bearer
// @Produce json
// @Param keyword query string false "关键词（工号/姓名/手机号）"
// @Param role query string false "角色过滤（employee/manager/admin）"
// @Param page query int false "页码，从1开始"
// @Param page_size query int false "每页数量，默认20，最大100"
// @Success 200 {object} utils.Response{data=dto.AdminListPointsResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/points [get]
func (h *PointHandler) AdminListAllPoints(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var query dto.AdminListPointsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	result, err := h.points.AdminListAllPoints(adminID, query)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	// 确保result不为nil
	if result == nil {
		result = &dto.AdminListPointsResponse{
			Items:      []dto.UserPointListItem{},
			Pagination: dto.Pagination{Page: query.Page, PageSize: query.PageSize, Total: 0},
		}
	}

	utils.NewSuccessResponse(result).JSON(c)
}
