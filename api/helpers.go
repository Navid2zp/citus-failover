package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func jsonBindErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error_code":   ErrCodeBadData,
		"message": "couldn't unpack the data you sent!",
		"error": err.Error(),
	})
}
