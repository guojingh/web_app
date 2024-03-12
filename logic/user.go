package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/snowflake"
)

// SignUp 存放业务逻辑处理代码
func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//生成uid	//密码加密
	userId := snowflake.GetID()
	//构造一个User实例
	user := &models.User{
		UserID:   userId,
		Username: p.Username,
		Password: p.Password,
	}

	//保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) error {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(user)
}
