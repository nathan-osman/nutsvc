package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func panicToJSONError(c *gin.Context, i interface{}) {
	var message string
	switch v := i.(type) {
	case error:
		message = v.Error()
	case string:
		message = v
	default:
		message = "an unknown error has occurred"
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}
