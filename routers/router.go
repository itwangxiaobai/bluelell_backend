package routers

import (
	"bluelell_backend/controller"
	"bluelell_backend/logger"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 请求注册路由
	r.POST("/signup", controller.SignUpHandler)
	return r
}
