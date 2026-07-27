package mysql

import (
	"errors"
	"mint/errs"
	"mint/models"

	"gorm.io/gorm"
)

func CreatePost(post *models.CreatePostParam, userID uint64, postID uint64) error {
	postModel := &models.Post{
		PostID:      postID,
		AuthorID:    userID,
		CommunityID: post.CommunityID,
		Status:      0,
		Title:       post.Title,
		Content: &models.PostContent{
			Content: post.Content,
		},
	}

	err := db.Create(postModel).Error

	if err != nil {
		return errs.NewAppError(500, errs.CodeInsertDBError, "将帖子数据插入数据库失败,事务已经回滚", err)
	}
	return nil
}

func GetPostDetailByID(postID uint64) (*models.Post, error) {
	var post models.Post

	err := db.Model(&models.Post{}).Joins("Content").Joins("Author").Joins("MainCommunity").Where("post.post_id = ?", postID).First(&post).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewAppError(404, errs.CodePostRecordNotFound, "对应ID的帖子不存在", nil)
		}

		return nil, errs.NewAppError(500, errs.CodeDBError, "数据库查询失败", err)
	}

	return &post, nil
}
