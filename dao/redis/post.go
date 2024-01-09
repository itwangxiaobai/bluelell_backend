package redis

import (
	"bluelell_backend/models"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis中获取id
	// 1.根据用户的请求参数中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2.确定查询的索引起始点
	return getIDsFormKey(key, p.Page, p.Size)

}

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(context.Background(), key, start, end).Result()
}

func GetPostVoteData(ids []string) (data []int64, err error) {
	/*data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		// 查找key中分数是1的元素的数量--->统计每篇帖子的赞成票的数量
		v := client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}*/
	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(context.Background(), key, "1", "1")
	}
	cmders, err := pipeline.Exec(context.Background())
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 根据社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用 zinterstore 把分区的帖子set于帖子分数的zset生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据
	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少zinterstroe执行的数据
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(context.Background(), key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(context.Background(), key, &redis.ZStore{
			Keys:      []string{cKey, orderKey},
			Aggregate: "MAX",
		})
		pipeline.Expire(context.Background(), key, 60*time.Second)
		_, err := pipeline.Exec(context.Background())
		if err != nil {
			return nil, err
		}
	}
	// 存在的话就直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
