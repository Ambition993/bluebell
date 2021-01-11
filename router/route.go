package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_app_base/controller"
	"web_app_base/logger"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//测试路由
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})
	//注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	// 登录业务路由
	r.POST("/signin", controller.SignInHandler)
	return r

}
