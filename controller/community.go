package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
)

// ---------和社区相关的

func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区 (community_id, community_name) 以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		//不轻易把服务端错误暴露给外面
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区详情分类
func CommunityDetailHandler(c *gin.Context) {
	//1.获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetailByID(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		//不轻易把服务端错误暴露给外面
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
