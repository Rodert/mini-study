package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/javapub/mini-study/mini-study-backend/internal/service"
	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

// UploadHandler handles file upload operations.
type UploadHandler struct {
	audit     *service.AuditService
	uploadDir string
	maxSize   int64
}

// NewUploadHandler creates a handler.
func NewUploadHandler(uploadDir string, maxSizeMB int, audit *service.AuditService) *UploadHandler {
	return &UploadHandler{audit: audit, uploadDir: uploadDir, maxSize: int64(maxSizeMB) * 1024 * 1024}
}

// Upload godoc
// @Summary 上传文件
// @Description 登录用户上传文件后返回可用于内容配置的存储路径
// @Tags 文件
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "待上传文件"
// @Success 200 {object} utils.Response{data=map[string]string}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/files/upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "file is required").JSON(c)
		return
	}

	if file.Size > h.maxSize {
		utils.NewErrorResponse(http.StatusBadRequest, "file too large").JSON(c)
		return
	}

	path, err := utils.SaveUploadedFile(file, h.uploadDir)
	if err != nil {
		utils.NewErrorResponse(http.StatusInternalServerError, err.Error()).JSON(c)
		return
	}

	_ = h.audit.Record(0, "upload", "files", file.Filename, "success")
	utils.NewSuccessResponse(gin.H{"path": path}).JSON(c)
}
