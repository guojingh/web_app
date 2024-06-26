package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 注册参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登陆参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	//UserID可以从当前登陆的用户获取
	PostID    string `json:"post_id" binding:"required"`              //帖子ID
	Direction int8   `json:"direction,string" binding:"oneof=-1 0 1"` //赞成票（1）还是反对票（-1）取消投票（0）
}

// ParamPostList 获取帖子列表 query string 参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"` //可以为空
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}

// ParamCommunityPostList 按社区获取帖子列表 query string 参数
type ParamCommunityPostList struct {
	*ParamPostList
}
