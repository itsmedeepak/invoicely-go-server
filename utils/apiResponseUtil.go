package utils

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ApiResponse(c *gin.Context, status int, success bool, message string, data interface{}) {
	c.JSON(status, bson.M{
		"success": success,
		"message": message,
		"data":    data,
	})
}
