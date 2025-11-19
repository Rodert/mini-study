package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

// Recovery captures panics and converts them into JSON responses.
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", zap.Any("error", r))
				utils.NewErrorResponse(http.StatusInternalServerError, "internal server error").JSON(c)
				c.Abort()
			}
		}()
		c.Next()
	}
}
