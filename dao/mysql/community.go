package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"bluebell/models"
)

// GetCommunityList 获取所有的community id 和 name
func GetCommunityList() (communityList []*models.Community, err error) {
	//在数据库里面找查询所有的community信息并且返回给logic层
	sqlStr := `select community_id , community_name from community`
	if err := db.Select(&communityList, sqlStr); err != nil {
		zap.L().Error("there is some err when selecting  in database", zap.Error(err))
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in database")
			err = nil
		}
	}
	return communityList, err
}

//GetCommunityDetailByID 根据ID 查询community的详细信息
func GetCommunityDetailByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	sqlStr := `select id , community_name , introduction, create_time from community where community_id = ?`
	communityDetail = new(models.CommunityDetail)
	if err := db.Get(communityDetail, sqlStr, id); err != nil {
		zap.L().Error("there is some err when selecting  in database", zap.Error(err))
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in database")
			err = ErrorInvalidID
		}
	}
	return
}
