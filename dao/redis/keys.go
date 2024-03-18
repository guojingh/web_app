package redis

// KeyPostTimeZSet redis key尽量使用命名空间的方式，方便业务拆分
const (
	KeyPrefix           = "bluebull:"
	KeyPostTimeZSet     = "post:time"   //zset；帖子以发帖时间为分数
	KeyPostScoreZSet    = "post:score"  //zset；帖子及投票的分数
	KeyPostVotedZSetPre = "post:voted:" //zset；记录用户及投票类型；参数是post_id
	KeyCommunitySetPre  = "community:"  //set;保存每个分区帖子下的id
)

// Redis Key 加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
