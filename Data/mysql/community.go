package mysql

import (
	"mint/errs"
	"mint/models"

	"gorm.io/gorm"
)

func GetCommunityList() ([]models.CommunityInfo, error) {
	var communityList []models.CommunityInfo
	result := db.Model(&models.CommunityInfo{}).Find(&communityList)
	if result.Error != nil {
		return nil, &errs.AppError{
			HTTPCode: 500,
			BizCode:  errs.CodeDBError,
			Message:  "查询数据库失败",
			Err:      result.Error,
		}
	}
	return communityList, nil
}

func GetCommunityDetailByID(communityID uint64) (*models.CommunityDetail, error) {
	var communityDetail models.CommunityDetail
	result := db.Model(&models.CommunityDetail{}).Where("community_id = ?", communityID).First(&communityDetail)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &errs.AppError{
				HTTPCode: 404,
				BizCode:  errs.CodeCommunityRecordNotFound,
				Message:  "ID对应的社区不存在",
				Err:      result.Error,
			}
		} else {
			return nil, &errs.AppError{
				HTTPCode: 500,
				BizCode:  errs.CodeDBError,
				Message:  "查询数据库失败",
				Err:      result.Error,
			}
		}

	}
	return &communityDetail, nil
}

func FindCommunityExists(communityID uint64) error {
	var dummy int
	err := db.Table("community").Select("1").Where("community_id = ?", communityID).Take(&dummy).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errs.AppError{
				HTTPCode: 404,
				BizCode:  errs.CodeCommunityRecordNotFound,
				Message:  "ID对应的社区不存在",
				Err:      nil,
			}
		} else {
			return &errs.AppError{
				HTTPCode: 500,
				BizCode:  errs.CodeDBError,
				Message:  "查询数据库失败",
				Err:      err,
			}
		}
	}
	return nil
}
