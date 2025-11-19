package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS configures cross-origin resource sharing rules.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	if len(allowedOrigins) == 0 {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = allowedOrigins
	}
	return cors.New(cfg)
}
