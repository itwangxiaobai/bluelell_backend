package redis

// redis key 注意命名空间的方式，方便查询和拆分

const (
	Prefix             = "bluebell:"   // 项目key的前缀
	KeyPostTimeZSet    = "post:time"   // zset:帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  // zset:帖子及帖子分数
	KeyPostVotedZSetPF = "post:voted:" // zset:记录用户及投票类型；参数是post_id
	KeyCommunitySetPF  = "community:"  // set:保存每个分区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}
