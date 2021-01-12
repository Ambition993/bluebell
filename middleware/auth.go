package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"web_app_base/pkg/jwt"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//获取authHeader
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "authHeader is empty",
			})
			c.Abort()
			return
		}
		// 切分authHeader 获得token那一段
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "authHeader 格式错误",
			})
			c.Abort()
			return
		}
		//解析token  并把username绑定到上下文中
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "token is invalid ",
			})
			c.Abort()
			return
		}
		c.Set("username", mc.Username)
		c.Next()
	}
}
