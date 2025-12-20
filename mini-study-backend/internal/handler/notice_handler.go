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

// NoticeHandler handles notice endpoints.
type NoticeHandler struct {
	service *service.NoticeService
}

// NewNoticeHandler creates handler.
func NewNoticeHandler(service *service.NoticeService) *NoticeHandler {
	return &NoticeHandler{service: service}
}

// GetLatestNotice godoc
// @Summary 获取最新系统公告
// @Description 返回当前时间范围内最新启用的公告（如有）
// @Tags 公告
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=dto.NoticeResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/notices/latest [get]
func (h *NoticeHandler) GetLatestNotice(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	notice, err := h.service.GetLatestNotice(userID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	if notice == nil {
		utils.NewSuccessResponse(nil).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toResponse(notice)).JSON(c)
}

// AdminListNotices godoc
// @Summary 管理员查询公告列表
// @Description 管理员可按状态筛选公告
// @Tags 管理后台-公告
// @Security Bearer
// @Produce json
// @Param status query bool false "是否启用 true/false"
// @Success 200 {object} utils.Response{data=[]dto.NoticeResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/notices [get]
func (h *NoticeHandler) AdminListNotices(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var query dto.AdminListNoticeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	notices, err := h.service.AdminListNotices(adminID, query.Status)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toResponses(notices)).JSON(c)
}

// AdminCreateNotice godoc
// @Summary 管理员创建公告
// @Description 管理员新增系统公告，支持文字和图片
// @Tags 管理后台-公告
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.AdminCreateNoticeRequest true "公告信息"
// @Success 200 {object} utils.Response{data=dto.NoticeResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/notices [post]
func (h *NoticeHandler) AdminCreateNotice(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var req dto.AdminCreateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	notice, err := h.service.AdminCreateNotice(adminID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toResponse(notice)).JSON(c)
}

// AdminUpdateNotice godoc
// @Summary 管理员更新公告
// @Description 管理员可修改公告内容、状态及生效时间
// @Tags 管理后台-公告
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Param body body dto.AdminUpdateNoticeRequest true "公告信息"
// @Success 200 {object} utils.Response{data=dto.NoticeResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/notices/{id} [put]
func (h *NoticeHandler) AdminUpdateNotice(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idParam := c.Param("id")
	noticeID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的公告ID").JSON(c)
		return
	}

	var req dto.AdminUpdateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	notice, err := h.service.AdminUpdateNotice(adminID, uint(noticeID), req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toResponse(notice)).JSON(c)
}

func (h *NoticeHandler) toResponses(items []model.Notice) []dto.NoticeResponse {
	resp := make([]dto.NoticeResponse, 0, len(items))
	for i := range items {
		resp = append(resp, h.toResponse(&items[i]))
	}
	return resp
}

func (h *NoticeHandler) toResponse(notice *model.Notice) dto.NoticeResponse {
	return dto.NoticeResponse{
		ID:        notice.ID,
		Title:     notice.Title,
		Content:   notice.Content,
		ImageURL:  notice.ImageURL,
		Status:    notice.Status,
		StartAt:   notice.StartAt,
		EndAt:     notice.EndAt,
		CreatedAt: notice.CreatedAt,
	}
}
