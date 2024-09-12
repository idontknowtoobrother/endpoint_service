package utils

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, errNum int, err error) {
	c.JSON(errNum, gin.H{
		"error": err.Error(),
	})
}

func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"data": data,
	})
}
