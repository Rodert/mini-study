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

// ExamHandler exposes exam endpoints.
type ExamHandler struct {
	service *service.ExamService
}

// NewExamHandler builds handler.
func NewExamHandler(svc *service.ExamService) *ExamHandler {
	return &ExamHandler{service: svc}
}

// ListAvailable godoc
// @Summary 获取考试列表
// @Tags 考试
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.ExamListItem}
// @Router /api/v1/exams [get]
func (h *ExamHandler) ListAvailable(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}
	resp, err := h.service.ListAvailableExams(userID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// GetExamDetail godoc
// @Summary 获取考试详情
// @Tags 考试
// @Security Bearer
// @Produce json
// @Param id path int true "考试ID"
// @Success 200 {object} utils.Response{data=dto.ExamDetailResponse}
// @Router /api/v1/exams/{id} [get]
func (h *ExamHandler) GetExamDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	examID, err := parseIDParam(c.Param("id"))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的考试ID").JSON(c)
		return
	}

	resp, err := h.service.GetExamDetail(userID, examID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// SubmitExam godoc
// @Summary 提交考试作答
// @Tags 考试
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "考试ID"
// @Param body body dto.ExamSubmitRequest true "提交答案"
// @Success 200 {object} utils.Response{data=dto.ExamSubmitResponse}
// @Router /api/v1/exams/{id}/submit [post]
func (h *ExamHandler) SubmitExam(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	examID, err := parseIDParam(c.Param("id"))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的考试ID").JSON(c)
		return
	}

	var req dto.ExamSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp, err := h.service.SubmitExam(userID, examID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// ListMyResults godoc
// @Summary 查询个人考试成绩
// @Tags 考试
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.ExamResultSummary}
// @Router /api/v1/exams/my/results [get]
func (h *ExamHandler) ListMyResults(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	resp, err := h.service.ListMyResults(userID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// AdminCreateExam godoc
// @Summary 管理员创建考试
// @Tags 管理端/考试
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.AdminExamUpsertRequest true "考试"
// @Success 200 {object} utils.Response{data=dto.ExamDetailResponse}
// @Router /api/v1/admin/exams [post]
func (h *ExamHandler) AdminCreateExam(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var req dto.AdminExamUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp, err := h.service.AdminCreateExam(adminID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// AdminListExams godoc
// @Summary 管理员获取考试列表
// @Tags 管理端/考试
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.ExamDetailResponse}
// @Router /api/v1/admin/exams [get]
func (h *ExamHandler) AdminListExams(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	resp, err := h.service.AdminListExams(adminID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// AdminGetExam godoc
// @Summary 管理员获取考试详情
// @Tags 管理端/考试
// @Security Bearer
// @Produce json
// @Param id path int true "考试ID"
// @Success 200 {object} utils.Response{data=dto.ExamDetailResponse}
// @Router /api/v1/admin/exams/{id} [get]
func (h *ExamHandler) AdminGetExam(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	examID, err := parseIDParam(c.Param("id"))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的考试ID").JSON(c)
		return
	}

	resp, err := h.service.AdminGetExam(adminID, examID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// AdminUpdateExam godoc
// @Summary 管理员更新考试
// @Tags 管理端/考试
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "考试ID"
// @Param body body dto.AdminExamUpsertRequest true "考试"
// @Success 200 {object} utils.Response{data=dto.ExamDetailResponse}
// @Router /api/v1/admin/exams/{id} [put]
func (h *ExamHandler) AdminUpdateExam(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	examID, err := parseIDParam(c.Param("id"))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的考试ID").JSON(c)
		return
	}

	var req dto.AdminExamUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp, err := h.service.AdminUpdateExam(adminID, examID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// ManagerOverview godoc
// @Summary 店长查看员工考试与学习进度
// @Tags 店长
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=dto.ManagerExamOverviewResponse}
// @Router /api/v1/manager/exams/overview [get]
func (h *ExamHandler) ManagerOverview(c *gin.Context) {
	managerID := middleware.GetUserID(c)
	if managerID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	resp, err := h.service.GetManagerOverview(managerID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

func parseIDParam(raw string) (uint, error) {
	id64, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

