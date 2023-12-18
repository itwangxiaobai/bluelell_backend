package logic

import (
	"bluelell_backend/dao/mysql"
	"bluelell_backend/models"
	"bluelell_backend/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) {
	// 1、判断用户存不存在
	mysql.QueryUserByUsername()
	// 2、生成UID
	snowflake.GetID()
	// 3、保存进数据库
	mysql.InsertUser()
}
