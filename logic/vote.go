package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"go.uber.org/zap"
	"strconv"
)

// PostVote 投票功能的实现

// 用户投票
// 使用简化版的投票算法
// 投一票就加432分    86400 / 200  200 个赞成了可以置顶一天了

/*
direction = 1
	1之前没投过 现在投赞成     更新分数和投票记录   +432       差值为1
	2之前投反对票 现在投赞成   更新分数和投票记录	+432 * 2   差值为2
direction = 0
	1之前投赞成 现在要取消     更新分数和投票记录   -432    差值为-1
	2之前投反对 现在要取消     更新分数和投票记录   +432    差值为1
direction = -1
	1之前没投过 现在反对      更新分数和投票记录    -432      差值为-1
	2之前赞成票 现在反对      更新分数和投票记录    -432 * 2  差值的-2

投票的限制:
	帖子自发表之日起一个星期可以让用户投票 超过一个星期就不能再投票
		1 到期之后将redis里面保存的赞成票数和反对票数存到MySQL
		2 到期之后删除那个KeyPostVoteZSetPF

*/
//1, 用户投票的数据 为帖子投票的函数
func PostVote(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug(" PostVote(userID int64, p models.ParamVoteData)", zap.Int64("userID", userID), zap.Int64("postID", p.PostID), zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.FormatInt(userID, 10), strconv.FormatInt(p.PostID, 10), float64(p.Direction))
}
