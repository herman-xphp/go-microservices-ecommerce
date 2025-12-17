package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func ResponseSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ResponseError(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}
