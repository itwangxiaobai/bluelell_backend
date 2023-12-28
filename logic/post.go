package logic

import (
	"bluelell_backend/dao/mysql"
	"bluelell_backend/models"
	"bluelell_backend/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	// 1.生成post_id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	return mysql.CreatePost(p)
	// 3.返回
}
