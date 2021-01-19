package redis

//redis key
//redis 使用命名空间的方式区分 方便查询和拆分

const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"   //帖子以及发布时间
	KeyPostScoreZSet       = "post:score"  // 帖子以及帖子分数
	KeyPostVotedZSetPrefix = "post:voted:" // 记录用户以及投票的类型 参数是post id
)
// getRedisKey 给redis加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
