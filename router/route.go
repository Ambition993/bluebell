package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/vi")
	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登录业务路由
	v1.POST("/signin", controller.SignInHandler)
	//
	//应用中间件
	v1.Use(middleware.JWTAuthMiddleware())
	//测试路由
	v1.GET("/ping", middleware.JWTAuthMiddleware(), func(context *gin.Context) {
		// 登录的用户才能访问 判断请求头有没有有效的jwt
		//isLogin := true
		//authHeader := context.Request.Header.Get("Authorization")
		//if
		context.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		v1.POST("/vote", controller.PostVoteController)
	}
	return r

}
