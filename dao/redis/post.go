package redis

import (
	"strconv"
	"time"
	"web_app/models"

	"github.com/go-redis/redis"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {

	start := (page - 1) * size
	end := size - 1

	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//1.根据用户请求参数中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.确定查询的索引起始点
	//3.ZRevRange按分数从大到小的顺序查询指定数量的元素
	return getIDsFromKey(key, p.Page, p.Size)

}

// GetPostVoteData 获取帖子的投票量
func GetPostVoteData(ids []string) (data []int64, err error) {
	/*	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPre + id)
		//查找key中分数是1的元素的数量，统计每篇帖子的赞成数量
		v := client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}*/

	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPre + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmdrs, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmdrs))
	for _, cmdr := range cmdrs {
		v := cmdr.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区根据IDS查找数据
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//使用zinterstore把分区的贴子与帖子分数的zset生成一个新的zset
	//针对新的zset按之前的逻辑取数据

	//社区的key orderKey string, communityID, page, size int64
	cKey := getRedisKey(KeyCommunitySetPre + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在的话直接直接根据key查询redis
	return getIDsFromKey(key, p.Page, p.Size)
}
