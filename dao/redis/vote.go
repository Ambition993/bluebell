package redis

import (
	"errors"
	redis "github.com/go-redis/redis/v8"
	"math"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("投票重复")
)

func VoteForPost(userID, postID string, value float64) error {
	//1 判断投票限制
	//去redis取得帖子发布的时间
	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2 更新帖子分数
	// 先查询之前的投票记录
	ovalue := rdb.ZScore(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()

	// 如果这次投票的情况和之前的一样  比如之前赞成 现在也投赞成就不行
	if ovalue == value {
		return ErrVoteRepeated
	}
	// 现在 - 之前 的投票记录（正反） 得到增值方向  乘上增值的绝对值 就是增加(减少的值)
	var op float64
	if value > ovalue {
		op = 1
	} else {
		op = -1
	}
	//计算两次投票的差值

	pipeline := rdb.TxPipeline()
	diff := math.Abs(ovalue - value)
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	//3 记录该用户为该帖子投票的数

	rz := redis.Z{Member: userID, Score: value}
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), &rz)
	}
	_, err := pipeline.Exec(ctx)
	return err
}
