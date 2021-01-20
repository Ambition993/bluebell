package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	//生成PostID
	p.ID = snowflake.GenID()
	// 保存到数据里面
	err = mysql.CreatePost(p)
	// 在redis里面记录帖子的创建时间
	err = redis.CreatePost(p.ID, p.CommunityID)
	//返回
	return err
}

// GetPostDetailByID 根据ID获取post详情
func GetPostDetailByID(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并且组合我们需要的数据
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Error(err))
	}
	// 根据查询帖子返回的UserID 查询User的详细信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
	}
	// 根据查询帖子返回的CommunityID 查询CommunityID的详细信息
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID)",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}
	return
}

// GetPostList 分页获取帖子
func GetPostList(page int64, size int64) (data []*models.ApiPostDetail, err error) {
	// 获取帖子列表
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	// 组合需要的数据
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据查询帖子返回的UserID 查询User的详细信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
		}
		// 根据查询帖子返回的CommunityID 查询CommunityID的详细信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID)",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 获取帖子 通过不同的排序方式   （byTime, byScore)
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis里面找id列表
	IDs, err := redis.GetPostIDsInorder(p)
	if err != nil {
		return
	}
	if len(IDs) == 0 {
		zap.L().Warn(" redis.GetPostIDsInorder(p) success but IDs' length is 0  return nothing")
		return
	}
	// 根据ID去MySQL查询帖子详情
	// 返回的数据还要根据postid来排序
	posts, err := mysql.GetPostListByIDs(IDs)
	// 封装详细post详情信息

	// 提前查询好redis里面的投票数量
	voteData, err := redis.GetPostVoteData(IDs)
	if err != nil {
		return
	}
	for idx, post := range posts {
		// 根据查询帖子返回的UserID 查询User的详细信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
		}
		// 根据查询帖子返回的CommunityID 查询CommunityID的详细信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID)",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
			Votes:           voteData[idx],
		}
		data = append(data, postDetail)
	}
	return
}

// GetCommunityPostList 通过选择社区获得帖子（byTime, byScore)
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis里面找id列表
	IDs, err := redis.GetCommunityPostIDsInorder(p)
	if err != nil {
		return
	}
	if len(IDs) == 0 {
		zap.L().Warn(" redis.GetPostIDsInorder(p) success but IDs' length is 0  return nothing")
		return
	}
	// 根据ID去MySQL查询帖子详情
	// 返回的数据还要根据postid来排序
	posts, err := mysql.GetPostListByIDs(IDs)
	// 封装详细post详情信息

	// 提前查询好redis里面的投票数量
	voteData, err := redis.GetPostVoteData(IDs)
	if err != nil {
		return
	}
	for idx, post := range posts {
		// 根据查询帖子返回的UserID 查询User的详细信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
		}
		// 根据查询帖子返回的CommunityID 查询CommunityID的详细信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID)",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
			Votes:           voteData[idx],
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListBy 将两个查询逻辑合二为一的函数  根据communityID的存在 来执行不同的查询逻辑
func GetPostListBy(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListBy() failed", zap.Error(err))
		return nil, err
	}
	return data, err
}
