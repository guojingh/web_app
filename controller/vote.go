package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"web_app/logic"
	"web_app/models"
)

type VoteData struct {
	//UserID可以从当前登陆的用户获取
	PostID    int64 `json:"post_id,string"`   //帖子ID
	Direction int   `json:"direction,string"` //赞成票（1）还是反对票（-1）
}

// PostVoteController 投票
func PostVoteController(c *gin.Context) {
	//参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			zap.L().Debug("PostVoteController failed", zap.Error(err))
			ResponseError(c, CodeInvalidParam)
			return
		}
		//翻译并去除掉错误信息中的结构体标识
		errData := removeTopStruct(errs.Translate(trans))
		zap.L().Debug("PostVoteController failed")
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//具体投票的业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
