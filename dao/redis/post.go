package redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)

func CreatePost(postID int64) error {
	// 创建帖子的时候要初始化帖子的创建时间以及初始分数
	// 用事务来操作
	pipeline := rdb.TxPipeline()
	//  帖子时间
	zTime := redis.Z{Score: float64(time.Now().Unix()), Member: postID}
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &zTime).Result()
	// 帖子分数
	zScore := redis.Z{Score: float64(time.Now().Unix()), Member: postID}
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), &zScore).Result()
	_, err := pipeline.Exec(ctx)
	return err
}
