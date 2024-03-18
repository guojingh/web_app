package mysql

import (
	"testing"
	"web_app/models"
	"web_app/settings"
)

func init() {
	dbCfg := settings.Mysql{
		Host:         "192.168.222.131",
		Port:         3306,
		User:         "root",
		Password:     "123456",
		DBName:       "bluebull",
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
		AuthorID:    123,
		CommunityID: 1,
		Status:      1,
		Title:       "test",
		Content:     "just a test",
	}

	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("createPost insert record into mysql failed, err:%v\n", err)
	}

	t.Logf("createPost insert record into mysql success")
}
