package logic

import (
	"bluelell_backend/dao/mysql"
	"bluelell_backend/models"
	"bluelell_backend/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1、判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2、生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID:   userID,
		Password: p.Password,
		Username: p.Username,
	}
	// 3、保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) error {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(user)
}
