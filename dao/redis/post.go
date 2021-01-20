package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

//CreatePost 创建帖子的时候要初始化帖子的创建时间以及初始分数
func CreatePost(postID, communityID int64) error {
	// 创建帖子的时候要初始化帖子的创建时间以及初始分数
	// 用事务来操作
	pipeline := rdb.TxPipeline()
	//  帖子时间
	zTime := redis.Z{Score: float64(time.Now().Unix()), Member: postID}
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &zTime).Result()
	// 帖子分数
	zScore := redis.Z{Score: float64(time.Now().Unix()), Member: postID}
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), &zScore).Result()
	// 补充：把帖子的ID添加到社区的set里面去
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ctx, cKey, postID)
	_, err := pipeline.Exec(ctx)
	return err
}

// getIDsFromKey 一个通过key 在redis zset里面取值的一个通用方法  会被以下的方法调用
func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 查询 从大到小返回数据
	return rdb.ZRevRange(ctx, key, start, end).Result()
}

// GetPostIDsInorder 从redis里面获取id
func GetPostIDsInorder(p *models.ParamPostList) ([]string, error) {
	// 从redis里面获取id
	// 根据用户传来的参数order来确定要查询的 redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderByScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 确定查询的索引起始点
	// 返回数据
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids获取每篇帖子的投赞成票的票数
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	pipeline := rdb.TxPipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInorder 按社区来查询ids
func GetCommunityPostIDsInorder(p *models.ParamPostList) (ids []string, err error) {
	// 使用分区zinterstore 吧分区的帖子和帖子分数的zset 生成一个新的zset
	// 针对新的zset按照之前的逻辑取得数据
	// 利用缓存Key减少zinterstore执行的次数
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderByScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 获取社区的id
	pipeline := rdb.TxPipeline()
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	// 利用缓存Key减少zinterstore执行的次数

	// 这个key就是需要 一个set 一个zset聚合起来的新的表的key
	key := orderKey + strconv.Itoa(int(p.CommunityID))

	if rdb.Exists(ctx, key).Val() < 1 {
		// 不存在 需要计算
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Aggregate: "MAX",
			Keys:      []string{cKey, orderKey},
		})                                        // zinterstore 计算
		pipeline.Expire(ctx, key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	// 存在就根据ID查询
	return getIDsFromKey(key, p.Page, p.Size)
}
