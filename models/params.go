package models

const (
	OrderByTime  = "time"
	OrderByScore = "score"
)

// 定义请求参数的结构体
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}
type ParamSignIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ParamVoteData struct {
	//可以在请求里面获取UseID
	PostID    int64 `json:"post_id" binding:"required"`       //postID
	Direction int8  `json:"direction" binding:"oneof=0 1 -1"` // 1 赞成 -1 反对
}

//获取帖子列表的一些参数
type ParamPostList struct {
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
	CommunityID int64  `json:"community_id" form:"community_id"` //可以为空
}

// ParamCommunityPostList 按社区获取帖子列表的query string参数
type ParamCommunityPostList struct {
	ParamPostList
}
