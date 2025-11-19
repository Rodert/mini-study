package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/middleware"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/service"
	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

// BannerHandler handles banner endpoints.
type BannerHandler struct {
	service *service.BannerService
}

// NewBannerHandler creates handler.
func NewBannerHandler(service *service.BannerService) *BannerHandler {
	return &BannerHandler{service: service}
}

// ListVisibleBanners godoc
// @Summary 查询可见轮播图
// @Description 返回当前登录用户基于角色可见的轮播图列表
// @Tags 轮播
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.BannerResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/banners [get]
func (h *BannerHandler) ListVisibleBanners(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	banners, err := h.service.ListVisible(userID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toResponses(banners)).JSON(c)
}

// AdminListBanners godoc
// @Summary 管理员查询轮播图列表
// @Description 管理员可按状态筛选轮播图
// @Tags 管理后台-轮播
// @Security Bearer
// @Produce json
// @Param status query bool false "是否启用 true/false"
// @Success 200 {object} utils.Response{data=[]dto.BannerResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/banners [get]
func (h *BannerHandler) AdminListBanners(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var query dto.AdminListBannerQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	banners, err := h.service.AdminListBanners(adminID, query.Status)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toResponses(banners)).JSON(c)
}

// AdminCreateBanner godoc
// @Summary 管理员创建轮播图
// @Description 管理员上传图片并配置跳转链接、可见角色等信息
// @Tags 管理后台-轮播
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.AdminCreateBannerRequest true "轮播图信息"
// @Success 200 {object} utils.Response{data=dto.BannerResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/banners [post]
func (h *BannerHandler) AdminCreateBanner(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var req dto.AdminCreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	banner, err := h.service.AdminCreateBanner(adminID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toResponse(banner)).JSON(c)
}

// AdminUpdateBanner godoc
// @Summary 管理员更新轮播图
// @Description 管理员可修改轮播图基础信息、状态及可见范围
// @Tags 管理后台-轮播
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "轮播图ID"
// @Param body body dto.AdminUpdateBannerRequest true "轮播图信息"
// @Success 200 {object} utils.Response{data=dto.BannerResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/banners/{id} [put]
func (h *BannerHandler) AdminUpdateBanner(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idParam := c.Param("id")
	bannerID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的轮播图ID").JSON(c)
		return
	}

	var req dto.AdminUpdateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	banner, err := h.service.AdminUpdateBanner(adminID, uint(bannerID), req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toResponse(banner)).JSON(c)
}

func (h *BannerHandler) toResponses(items []model.Banner) []dto.BannerResponse {
	resp := make([]dto.BannerResponse, 0, len(items))
	for i := range items {
		resp = append(resp, h.toResponse(&items[i]))
	}
	return resp
}

func (h *BannerHandler) toResponse(banner *model.Banner) dto.BannerResponse {
	return dto.BannerResponse{
		ID:           banner.ID,
		Title:        banner.Title,
		ImageURL:     banner.ImageURL,
		LinkURL:      banner.LinkURL,
		VisibleRoles: banner.VisibleRoles,
		SortOrder:    banner.SortOrder,
		Status:       banner.Status,
		StartAt:      banner.StartAt,
		EndAt:        banner.EndAt,
	}
}
