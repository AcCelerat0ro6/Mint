package logic

import (
	"mint/Data/mysql"
	"mint/models"
)

func GetCommunityList() ([]models.CommunityInfo, error) {
	return mysql.GetCommunityList()
}

func GetCommunityByID(communityID uint64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(communityID)
}
