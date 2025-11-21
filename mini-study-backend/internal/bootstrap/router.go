package bootstrap

import (
	"github.com/gin-gonic/gin"

	"github.com/javapub/mini-study/mini-study-backend/internal/handler"
	"github.com/javapub/mini-study/mini-study-backend/internal/middleware"
	"github.com/javapub/mini-study/mini-study-backend/internal/router"
)

// RegisterRoutes binds all HTTP handlers to the gin engine.
func RegisterRoutes(engine *gin.Engine, cfg *Config, userHandler *handler.UserHandler, contentHandler *handler.ContentHandler, learningHandler *handler.LearningHandler, bannerHandler *handler.BannerHandler, examHandler *handler.ExamHandler, uploadHandler *handler.UploadHandler, systemHandler *handler.SystemHandler, pointHandler *handler.PointHandler) {
	engine.Static("/uploads", cfg.Upload.Dir)
	auth := middleware.JWT(cfg.JWT.Secret)
	router.RegisterRoutes(engine, cfg.Swagger.Enabled, auth, userHandler, contentHandler, learningHandler, bannerHandler, examHandler, uploadHandler, systemHandler, pointHandler)
}
