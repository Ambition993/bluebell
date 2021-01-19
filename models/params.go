package models

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
	PostID    int64 `json:"post_id" binding:"required"`                 //postID
	Direction int8  `json:"direction" binding:"oneof=0 1 -1"` // 1 赞成 -1 反对
}
