package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_app_base/logger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"hello": "hello",
		})
	})
	return r
}
