// Package response
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    code,
		Message: message,
		Error:   message,
	})
}

func ServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Code:    50000,
		Message: "Internal server error",
		Error:   err.Error(),
	})
}
