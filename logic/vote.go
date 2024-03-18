package logic

import (
	"go.uber.org/zap"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"
)

// VoteForPost 投票功能 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/ 本项目使用简化版的投票
// 投一票就加423分  86400/200 -->200张赞成票可以给贴子续一天 ---> 《redis实战》
/* 投票的几种情况
direction=1时，有两种情况
	1.之前没有投过票，现在投赞成票  --> 更新分数和投票记录
	2.之前投反对票，现在改投赞成票  --> 更新分数和投票记录
direction=0时，有两种情况
	1.之前投赞成票，现在取消投票  --> 更新分数和投票记录
	2.之前投过反对票，现在取消投票  --> 更新分数和投票记录
direction=-1时，有两种情况
	1.之前没有投票，现在投反对票  --> 更新分数和投票记录
	2.之前投赞成票，现在改投反对票  --> 更新分数和投票记录

投票限制：
每个帖子发表之日起一个星期内允许用户投票，超过一个星期就不允许投票
	1.到期之后将redis中保存的赞成票数和反对票数存储到mysql中
	2.到期之后删除那个 KeyPostVotedZSetPre
*/
func VoteForPost(userID int64, p *models.ParamVoteData) error {

	zap.L().Debug("VoteForPost", zap.Int64("userID", userID), zap.String("postID", p.PostID), zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))

}
