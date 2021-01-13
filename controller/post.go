package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"bluebell/logic"
	"bluebell/models"
)

func CreatePostHandler(c *gin.Context) {
	// 1 参数校验
	p := new(models.Post)
	err := c.ShouldBindJSON(p) //validator binding
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从 c 取出上下文里面的userID

	// 2 创建帖子
	err = logic.CreatePost(p)
	if err != nil {
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3 返回success
	ResponseSuccess(c, nil)
}
