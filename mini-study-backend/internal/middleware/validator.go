package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const validatorKey = "validator"

// Validator injects a shared validator into the request context.
func Validator(v *validator.Validate) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(validatorKey, v)
		c.Next()
	}
}

// GetValidator retrieves the validator instance from context.
func GetValidator(c *gin.Context) *validator.Validate {
	if v, exists := c.Get(validatorKey); exists {
		if cast, ok := v.(*validator.Validate); ok {
			return cast
		}
	}
	return validator.New()
}
