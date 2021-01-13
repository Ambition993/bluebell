package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
)
const CtxUserIDKey string = "userID"
var ErrorUserNotLogin = errors.New("用户未登录")

//GetCurrentUser 获取当前登录的用户ID
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
