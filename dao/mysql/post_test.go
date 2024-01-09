package mysql

import (
	"bluelell_backend/models"
	"bluelell_backend/settings"
	"testing"
)

func init() {
	dbCfg := settings.MysqlConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "admin",
		DB:           "bluebell",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          10,
		CommunityID: 1,
		AuthorID:    123,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
