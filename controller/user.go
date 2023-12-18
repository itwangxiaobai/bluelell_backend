package controller

import (
	"bluelell_backend/logic"
	"bluelell_backend/models"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	// 1、获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("signup invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	fmt.Println(p)
	// 2、业务处理
	logic.SignUp(p)
	// 3、返回相应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success...",
	})
}
