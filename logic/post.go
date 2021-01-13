package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	//生成PostID
	p.ID = snowflake.GenID()
	// 保存到数据里面
	err = mysql.CreatePost(p)
	//返回
	return err
}
