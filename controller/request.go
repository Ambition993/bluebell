package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const CtxUserIDKey string = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUserID 获取当前登录的用户ID
func GetCurrentUserID(c *gin.Context) (err error, userID int64) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
	}
	return
}

// getPageInfo  从参数中解析分页信息
func getPageInfo(c *gin.Context) (int64, int64) {
	var (
		page int64
		size int64
		err  error
	)
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")
	page, err = strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
