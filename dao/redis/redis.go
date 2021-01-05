package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"web_app_base/settings"
)

var rdb *redis.Client
var ctx = context.Background()

func Init(redisConfig *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			redisConfig.Host,
			redisConfig.Port,
		),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	})
	_, err = rdb.Ping(ctx).Result()
	return err
}
func Close() {
	_ = rdb.Close()
}
