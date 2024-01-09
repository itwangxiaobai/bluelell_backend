package routers

import (
	"bluelell_backend/controller"
	"bluelell_backend/logger"
	"bluelell_backend/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 后端请求接口 /api/v1
	v1 := r.Group("/api/v1")
	// 请求注册路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登陆路由
	v1.POST("/login", controller.LoginHandler)

	// 根据时间或分数获取帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)
	v1.GET("/posts", controller.GetPostListHandler)
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/post/:id", controller.GetPostDetailHandler)

	v1.Use(middlewares.JWTAuthMiddleware(), middlewares.RateLimitMiddleware(2*time.Second, 1))
	{
		// v1.GET("community", controller.CommunityHandler)
		// v1.GET("community/:id", controller.CommunityDetailHandler)
		v1.POST("post", controller.CreatePostHandler)
		// v1.GET("post/:id", controller.GetPostDetailHandler)
		// v1.GET("/posts", controller.GetPostListHandler)
		v1.POST("/vote", controller.PostVoteController)
		// v1.GET("/posts2", controller.GetPostListHandler2)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
