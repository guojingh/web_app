package controller

import (
	"errors"
	"fmt"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {

	//1.参数校验和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		zap.L().Error("login.signup with invalid param", zap.Error(err))
		//
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//手动对请求参数进行详细的业务规则校验
	/*	if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
		//请求参数有误
		zap.L().Error("signup with invalid param")
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}*/
	fmt.Println(p)

	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("login.signup with invalid param", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}

		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		zap.L().Error("login with invalid param", zap.Error(err))
		//
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // id值大于 1<<53-1 int64类型的最大值1<<63-1
		"user_name": user.Username,
		"token":     user.Token,
	})
}
