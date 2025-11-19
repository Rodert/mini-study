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

// LearningHandler handles learning progress endpoints.
type LearningHandler struct {
	service *service.LearningService
}

// NewLearningHandler creates a learning handler.
func NewLearningHandler(service *service.LearningService) *LearningHandler {
	return &LearningHandler{service: service}
}

// UpdateProgress godoc
// @Summary 记录学习进度
// @Description 登录用户上报视频播放位置，系统会累计已学进度
// @Tags 学习
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.LearningProgressRequest true "学习进度"
// @Success 200 {object} utils.Response{data=dto.LearningProgressResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/learning [post]
func (h *LearningHandler) UpdateProgress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var req dto.LearningProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp, err := h.service.UpdateProgress(userID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// GetProgress godoc
// @Summary 查询指定内容的学习进度
// @Description 返回当前用户在某个内容上的进度详情
// @Tags 学习
// @Security Bearer
// @Produce json
// @Param content_id path int true "内容ID"
// @Success 200 {object} utils.Response{data=dto.LearningProgressResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/learning/{content_id} [get]
func (h *LearningHandler) GetProgress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	contentIDParam := c.Param("content_id")
	contentID, err := strconv.ParseUint(contentIDParam, 10, 64)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的内容ID").JSON(c)
		return
	}

	resp, err := h.service.GetProgress(userID, uint(contentID))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// ListProgress godoc
// @Summary 查询学习进度列表
// @Description 返回当前用户所有已记录的学习进度
// @Tags 学习
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.LearningProgressResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/learning [get]
func (h *LearningHandler) ListProgress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	resp, err := h.service.ListProgress(userID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}
