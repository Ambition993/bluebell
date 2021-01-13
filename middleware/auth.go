package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"bluebell/controller"
	"bluebell/pkg/jwt"
)

// JWTAuthMiddleware 处理登录后获取token的中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//获取authHeader
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 切分authHeader 获得token那一段
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		//解析token  并把username绑定到上下文中
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		//保存到上下文中 可以c.Get(CtxUserIDKey) 获得userID
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next()
	}
}
