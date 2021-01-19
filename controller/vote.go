package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PostVoteController 处理投票的方法

func PostVoteController(c *gin.Context) {
	//参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //断言类型
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))  //翻译除去错误提示中的结构体
		ResponseErrorWithMsg(c, CodeInvalidParam, errData) // 返回错误数据到前端
		return
	}
	err , userID:= GetCurrentUserID(c)
	if err != nil{
		ResponseError(c, CodeNeedLogin)
	}

	// 具体的投票逻辑
	logic.PostVote(userID, p)
	ResponseSuccess(c, nil)
}
