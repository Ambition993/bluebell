package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CreatePostHandler 創建一個帖子的處理函數
func CreatePostHandler(c *gin.Context) {
	// 1 参数校验
	p := new(models.Post)
	err := c.ShouldBindJSON(p) //validator binding
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从 c 取出上下文里面的userID
	err, userID := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	p.AuthorID = userID
	// 2 创建帖子
	err = logic.CreatePost(p)
	if err != nil {
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3 返回success
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 處理帖子詳情的一個處理函數
func GetPostDetailHandler(c *gin.Context) {
	// 1 获取参数
	pidStr := c.Param("id")

	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with error", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2 根据ID取出帖子的数据
	data, err := logic.GetPostDetailByID(pid)
	if err != nil {
		zap.L().Error(" logic.GetPostDetail(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3 返回响应
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	//解析参数
	page, size := getPageInfo(c)
	//获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error(" logic.GetPostList() ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回数据
	ResponseSuccess(c, data)
}

// GetPostListHandler2
/*
	 根据前端传来的参数 选择获取帖子的方式 （时间 分数）
	1 获取参数
	2 去redis里面查询id列表
	3 根据ID去MySQL里面查询
*/
func GetPostListHandler2(c *gin.Context) {
	// 获取请求参数
	// 先预设一点默认参数
	p := &models.ParamPostList{
		Page: 1, Size: 10, Order: models.OrderByTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取数据
	data, err := logic.GetPostListBy(p)
	if err != nil {
		zap.L().Error(" logic.GetPostList() ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回数据
	ResponseSuccess(c, data)
}
//
func GetCommunityPostListHandler(c *gin.Context) {
	// 初始化参数
	var p = new(models.ParamPostList)
	p.Page = 1
	p.Size = 10
	p.Order = models.OrderByTime

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取数据
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error(" logic.GetPostList() ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回数据
	ResponseSuccess(c, data)
}
