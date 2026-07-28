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

func GetPostList(page, size int) ([]models.Post, error) {
	var postList []models.Post

	// 1. 计算偏移量
	offset := (page - 1) * size

	// 2. 进行延迟关联优化 ，避免查询所有关联数据。

	// 2.1 构建子查询 (仅查询自增主键 id 以实现覆盖索引，防止回表)
	subQuery := db.Model(&models.Post{}).Select("id").Order("created_at desc").Limit(size).Offset(offset)

	// 2.2 执行查询 (使用 id 进行延迟关联)
	err := db.Model(&models.Post{}).
		Joins("INNER JOIN (?) AS t2 ON post.id = t2.id", subQuery).
		Joins("Author").
		Joins("MainCommunity").
		Joins("Content").
		Order("post.created_at desc").
		Find(&postList).Error

	if err != nil {
		return nil, errs.NewAppError(500, errs.CodeDBError, "数据库查询失败", err)
	}

	if len(postList) == 0 {
		return nil, errs.NewAppError(404, errs.CodePostListEmpty, "查询目标范围内帖子列表为空", nil)
	}

	return postList, nil
}
