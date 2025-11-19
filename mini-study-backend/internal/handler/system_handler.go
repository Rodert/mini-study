package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

// SystemHandler exposes health/version endpoints.
type SystemHandler struct {
	appName    string
	appVersion string
}

// NewSystemHandler builds handler.
func NewSystemHandler(appName, version string) *SystemHandler {
	return &SystemHandler{appName: appName, appVersion: version}
}

// Health godoc
// @Summary 健康检查
// @Description 返回服务运行状态与当前时间
// @Tags 系统
// @Produce json
// @Success 200 {object} utils.Response
// @Router /healthz [get]
func (h *SystemHandler) Health(c *gin.Context) {
	utils.NewSuccessResponse(gin.H{
		"status":  "ok",
		"time":    time.Now().UTC(),
		"service": h.appName,
	}).JSON(c)
}

// Version godoc
// @Summary 查询服务版本
// @Description 返回当前部署的应用名称与版本
// @Tags 系统
// @Produce json
// @Success 200 {object} utils.Response
// @Router /version [get]
func (h *SystemHandler) Version(c *gin.Context) {
	utils.NewSuccessResponse(gin.H{
		"service": h.appName,
		"version": h.appVersion,
	}).JSON(c)
}
