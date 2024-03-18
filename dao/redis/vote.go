package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

/* 投票的几种情况
direction=1时，有两种情况
	1.之前没有投过票，现在投赞成票  --> 更新分数和投票记录  差值的绝对值：1 +432
	2.之前投反对票，现在改投赞成票  --> 更新分数和投票记录  差值的绝对值：2 +432*2
direction=0时，有两种情况
	1.之前投赞成票，现在取消投票  --> 更新分数和投票记录  差值的绝对值：1 -432
	2.之前投过反对票，现在取消投票  --> 更新分数和投票记录  差值的绝对值：1 +432
direction=-1时，有两种情况
	1.之前没有投票，现在投反对票  --> 更新分数和投票记录  差值的绝对值：1  -432
	2.之前投赞成票，现在改投反对票  --> 更新分数和投票记录  差值的绝对值：2 -432*2

投票限制：
每个帖子发表之日起一个星期内允许用户投票，超过一个星期就不允许投票
	1.到期之后将redis中保存的赞成票数和反对票数存储到mysql中
	2.到期之后删除那个 KeyPostVotedZSetPre
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {

	//开启redis事务
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//贴子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPre + strconv.Itoa(int(communityID)))
	client.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票的限制
	//去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2和3 需要放到一个事务里卖弄
	//2.更新分数
	//先查当前用户给当前帖子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPre+postID), userID).Val()
	//如果这一次投票的值和之前保存的值一致就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	//3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPre+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPre+postID), redis.Z{
			Score:  value, //当前的用户是赞成票还是反对票
			Member: userID,
		})
	}

	_, err := pipeline.Exec()
	return err
}
