package routes

import (
	"net/http"
	"web_app/controller"
	"web_app/logger"
	"web_app/middlewares"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	//r.Use(logger.GinL	ogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(time.Second*2, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1 := r.Group("/api/v1")
	v1.Use(middlewares.Cors())
	// 注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)
	// 获取帖子列表
	v1.GET("/posts", controller.GetPostListHandler)

	// 使用JWT中间件
	v1.Use(middlewares.JWTAuthMiddleware())
	//贴子分类
	{
		// 获取社区
		v1.GET("/community", controller.CommunityHandler)
		// 社区详情分类
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		// 获取贴子详情
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		//根据时间或者分数获取帖子列表
		v1.GET("/posts2", controller.GetPostListHandler2)
		// 创建帖子
		v1.POST("/post", controller.CreatePostHandler)
		//投票
		v1.POST("/vote", controller.PostVoteController)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	pprof.Register(r)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})
	return r
}
