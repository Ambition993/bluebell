package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//生成PostID
	p.ID = snowflake.GenID()
	// 保存到数据里面
	err = mysql.CreatePost(p)
	// 在redis里面记录帖子的创建时间
	err = redis.CreatePost(p.ID)
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
