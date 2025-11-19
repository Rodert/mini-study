package router

import (
	"github.com/gin-gonic/gin"

	"github.com/javapub/mini-study/mini-study-backend/internal/handler"
)

// RegisterSystemRoutes registers health/version endpoints.
func RegisterSystemRoutes(engine *gin.Engine, systemHandler *handler.SystemHandler) {
	engine.GET("/healthz", systemHandler.Health)
	engine.GET("/version", systemHandler.Version)
}
