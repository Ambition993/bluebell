package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	//在数据库里面查找所有的community 返回
	return mysql.GetCommunityList()
}
func GetCommunityDetail(id int64) (communityDetail *models.CommunityDetail, err error) {
	//在数据库里面根据ID查找communityDetail 返回
	return mysql.GetCommunityDetailByID(id)
}
