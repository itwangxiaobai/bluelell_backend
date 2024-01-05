package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"math"
	"strconv"
	"time"
)

/*
投票的几种情况
direction = 1
1、之前没有投票，现在投赞成票  --> 更新分数和投票记录	差值的绝对值 1 +432
2、之前投反对票，现在投赞成票  --> 更新分数和投票记录	差值的绝对值 2 +432*2
direction = 0
1、之前投反对票，现在取消投票  --> 更新分数和投票记录	差值的绝对值 1 +432
2、之前投赞成票，现在取消投票  --> 更新分数和投票记录	差值的绝对值 1 -432
direction = -1
1、之前没有投票，现在投反对票  --> 更新分数和投票记录	差值的绝对值 1 -432
2、之前投赞成票，现在投反对票  --> 更新分数和投票记录	差值的绝对值 2 -432*2

投票的限制：
   每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
   	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
   	2. 到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	onWeekInSeconds = 7 * 24 * 3600
	scorePerVote    = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(context.Background(), getRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(context.Background(), getRedisKey(KeyPostScoreZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 更新：把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(context.Background(), cKey, postID)
	_, err := pipeline.Exec(context.Background())
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 1、判断帖子限制
	// 去redis取帖子的发帖时间
	postTime := client.ZScore(context.Background(), getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > onWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2、更新帖子分数
	// 先查当前用户给当前帖子的投票记录
	ov := client.ZScore(context.Background(), getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	zap.L().Debug("ov and value的值是", zap.Float64("value", value), zap.Float64("ov", ov))

	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	// 2、3需要放在一个事务中进行处理
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(context.Background(), getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3、记录用户该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(context.Background(), getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipeline.ZAdd(context.Background(), getRedisKey(KeyPostVotedZSetPF+postID), &redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec(context.Background())
	return err
}
