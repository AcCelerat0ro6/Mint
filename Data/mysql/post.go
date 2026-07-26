package mysql

import (
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
	}
	contentModel := &models.PostContent{
		PostID:  postID,
		Content: post.Content,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(postModel).Error; err != nil {
			return err
		}
		if err := tx.Create(contentModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errs.NewAppError(500, errs.CodeInsertDBError, "将帖子数据插入数据库失败,事务已经回滚", err)
	}
	return nil
}
