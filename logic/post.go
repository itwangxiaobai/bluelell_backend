package logic

import (
	"bluelell_backend/dao/mysql"
	"bluelell_backend/models"
	"bluelell_backend/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1.生成post_id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	return mysql.CreatePost(p)
	// 3.返回
}

func GetPostById(pid int64) (post *models.Post, err error) {
	post, err = mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	return
}
