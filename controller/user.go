package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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
			ResponseError(c, CodeInvalidParam)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	fmt.Println(p)
	// 手动对参数进行详细的业务规则校验
	// 业务处理
	if err := logic.SignUp(&p); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, CodeSuccess)
}

func SignInHandler(c *gin.Context) {
	var p models.ParamSignIn
	if err := c.BindJSON(&p); err != nil {
		//参数有误 直接返回
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断是不是validator里面的错误类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 手动检验业务
	user, err := logic.SignIn(&p)
	if err != nil {
		zap.L().Error("logic.Sign failed ", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		} else {
			ResponseErrorWithMsg(c, CodeServerBusy, err)
			return
		}

	}
	//返回响应

	ResponseSuccess(c, gin.H{
		"user_id":   user.UserID,
		"user_name": user.Username,
		"token":     user.Token,
	})
}
