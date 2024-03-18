package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePostHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST("/api/v1/post", CreatePostHandler)

	body := `{
		"community_id": 1,
		"title": "test",
		"content": "just a test"
		
	}`
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	//判断响应的内容是不是按预期返回了需要登陆的错误

	//方法1：判断响应的内容中是不是包含指定字符串
	assert.Contains(t, w.Body.String(), "需要登陆")
	//方法2：将响应的内容反序列到ResponseData，然后判断各字段与预期是否一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	}

	assert.Equal(t, res.Code, CodeNeedLogin)
}