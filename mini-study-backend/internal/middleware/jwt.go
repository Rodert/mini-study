package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

const userIDKey = "userID"

// JWT protects routes using bearer tokens.
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			utils.NewErrorResponse(http.StatusUnauthorized, "missing token").JSON(c)
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(secret, token)
		if err != nil {
			utils.NewErrorResponse(http.StatusUnauthorized, "invalid token").JSON(c)
			c.Abort()
			return
		}

		c.Set(userIDKey, claims.UserID)
		c.Next()
	}
}

// GetUserID fetches the authenticated user id from context.
func GetUserID(c *gin.Context) uint {
	value, exists := c.Get(userIDKey)
	if !exists {
		return 0
	}
	if id, ok := value.(uint); ok {
		return id
	}
	return 0
}
