package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
	"web_app/models"
	"web_app/settings"
)

// CreatePostHandler 创建贴子
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	} // validator --> binding tag

	//从 c 取到当前发请求的用户 id 值
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取贴子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数（帖子ID 从URL中获取）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.根据帖子ID获取帖子数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("get post detail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	page, size := getPageInfo(c)
	//1.获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//2.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版获取帖子列表，根据前端传来的参数动态的获取动态列表（1.按分数  2.按创建时间）
// 1.获取请求的query string参数
// 2.去redis查询id列表
// 3.根据ID去数据库查询贴子详细信息
func GetPostListHandler2(c *gin.Context) {
	// get 请求参数：/api/v1/posts2?page=1&size=10&score=time
	p := &models.ParamPostList{
		Page:  settings.Conf.App.Page,
		Size:  settings.Conf.App.Size,
		Order: models.OrderTime,
	}
	p.Page, p.Size = getPageInfo(c)
	//c.shouldBind()  根据请求的数据类型选择对应的方法获取数据
	//c.ShouldBindJSON()如果请求中携带的是json格式的数据，才能用这个方法获取到
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2() with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//1.获取数据
	data, err := logic.GetPostListNew(p) //更新合二为一
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//2.返回响应
	ResponseSuccess(c, data)
}

/*// GetCommunityPostListHandler 根据社区查询帖子列表
func GetCommunityPostListHandler(c *gin.Context) {
	// get 请求参数：/api/v1/posts2?page=1&size=10&score=time
	p := &models.ParamCommunityPostList{
		ParamPostList: &models.ParamPostList{
			Page:  settings.Conf.App.Page,
			Size:  settings.Conf.App.Size,
			Order: models.OrderTime,
		},
	}
	p.Page, p.Size = getPageInfo(c)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler() with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//c.shouldBind()  根据请求的数据类型选择对应的方法获取数据
	//c.ShouldBindJSON()如果请求中携带的是json格式的数据，才能用这个方法获取到
	//1.获取数据
	data, err := logic.GetCommunityPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetCommunityPostListHandler() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//2.返回响应
	ResponseSuccess(c, data)
}*/
