package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/javapub/mini-study/mini-study-backend/internal/middleware"
)

// RegisterMiddlewares wires all global middlewares.
func RegisterMiddlewares(engine *gin.Engine, cfg *Config, logger *zap.Logger, validate *validator.Validate) {
	engine.Use(
		middleware.RequestID(),
		middleware.RequestLogger(logger),
		middleware.Recovery(logger),
		middleware.CORS(cfg.Server.AllowedOrigins),
		middleware.Validator(validate),
	)
}
