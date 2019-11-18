package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InternalError(c *gin.Context, message string, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status": gin.H{
			"code": http.StatusInternalServerError,
			"desc": message,
		},
		"data": nil,
	})
	log.Error(fmt.Sprintf("ERROR %s %s %s. ", c.Request.Method, c.Request.RequestURI, message), err)
}

func Forbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status": gin.H{
			"code": http.StatusUnauthorized,
			"desc": "Unauthorized",
		},
		"data": nil,
	})
}

func NotFound(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"status": gin.H{
			"code": http.StatusNotFound,
			"desc": "Resource not found",
		},
		"data": nil,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status": gin.H{
			"code": http.StatusBadRequest,
			"desc": message,
		},
		"data": nil,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"code": 200,
			"desc": "OK",
		},
		"data": data,
	})
}
