package routers

import (
	"bluelell_backend/controller"
	"bluelell_backend/logger"
	"bluelell_backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/v1")
	// 请求注册路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登陆路由
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("community", controller.CommunityHandler)
		v1.GET("community/:id", controller.CommunityDetailHandler)
		v1.POST("post", controller.CreatePostHandler)
		v1.GET("post/:id", controller.GetPostDetailHandler)
	}

	return r
}
