package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API response envelope.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON writes the response to the http client.
func (r Response) JSON(c *gin.Context) {
	c.JSON(r.Code, r)
}

// NewSuccessResponse wraps success payloads.
func NewSuccessResponse(data interface{}) Response {
	return Response{Code: http.StatusOK, Message: "success", Data: data}
}

// NewErrorResponse wraps error payloads.
func NewErrorResponse(status int, message string) Response {
	return Response{Code: status, Message: message}
}
