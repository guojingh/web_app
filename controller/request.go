package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const ContextUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登陆")

// 可以获取当前登陆用户的id
func getCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIDKey)

	if !ok {
		err = ErrorUserNotLogin
		return
	}

	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func getPageInfo(c *gin.Context) (int64, int64) {
	//获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64
		size int64
		err  error
	)

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
