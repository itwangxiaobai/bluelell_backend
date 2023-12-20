package routers

import (
	"bluelell_backend/controller"
	"bluelell_backend/logger"
	"bluelell_backend/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 请求注册路由
	r.POST("/signup", controller.SignUpHandler)
	// 登陆路由
	r.POST("/login", controller.LoginHandler)
	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})
	return r
}
