package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
)

func GetCommunityList() ([]*models.Community, error) {

	//查找到所有的community 并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetailByID(id)
}
