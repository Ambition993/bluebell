package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	"web_app_base/logic"
	"web_app_base/models"
)

/*
	用户模块的controller 处理关于用户的请求
*/

//处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 处理参数
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		//参数有误 直接返回
		zap.L().Error("SignUp with invalid param", zap.Error(err))

		//判断是不是validator里面的错误类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		}
		// validator.ValidationErrors类型错误则进行翻译
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}
	fmt.Println(p)
	// 手动对参数进行详细的业务规则校验
	// 业务处理
	if err := logic.SignUp(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func SignInHandler(c *gin.Context) {
	var p models.ParamSignIn
	if err := c.BindJSON(&p); err != nil {
		//参数有误 直接返回
		zap.L().Error("SignUp with invalid param", zap.Error(err))

		//判断是不是validator里面的错误类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		}
		// validator.ValidationErrors类型错误则进行翻译
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}
	// 手动检验业务
	if err := logic.SignIn(&p); err != nil {
		zap.L().Error("logic.Sign failed ", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户名或者密码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
	})
}
