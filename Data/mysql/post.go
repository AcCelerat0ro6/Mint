package mysql

import (
	"errors"
	"mint/errs"
	"mint/models"
	"strconv"

	"gorm.io/gorm"
)

// CreatePost 创建帖子
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

// DeletePostByID 删除指定帖子及其内容，用于创建帖子失败时的补偿回滚
func DeletePostByID(postID uint64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("post_id = ?", postID).Delete(&models.PostContent{}).Error; err != nil {
			return errs.NewAppError(500, errs.CodeDBError, "补偿删除帖子内容失败", err)
		}

		result := tx.Where("post_id = ?", postID).Delete(&models.Post{})
		if result.Error != nil {
			return errs.NewAppError(500, errs.CodeDBError, "补偿删除帖子主记录失败", result.Error)
		}
		if result.RowsAffected == 0 {
			return errs.NewAppError(500, errs.CodeDBError, "补偿删除帖子主记录失败", gorm.ErrRecordNotFound)
		}

		return nil
	})
}

// GetPostDetailByID 根据ID获取帖子详情
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

// GetPostList 获取帖子列表
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

// GetPostListByIdList 根据ID列表查询帖子列表，返回的帖子按 Redis 传进来的顺序组装
func GetPostListByIdList(postIDs []string) ([]*models.Post, error) {

	length := len(postIDs)
	queryPostIds := make([]uint64, 0, length)
	postList := make([]*models.Post, 0, length)
	postIDMap := make(map[uint64]*models.Post, length)

	// string 转换为 uint64 类型
	for _, postID := range postIDs {
		queryPostID, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return nil, errs.NewAppError(500, errs.CodeInternalError, "Redis返回的帖子ID格式错误", err)
		}
		queryPostIds = append(queryPostIds, queryPostID)
	}

	// 执行查询
	var unorderPostList []*models.Post
	err := db.Model(&models.Post{}).
		Where("post.post_id IN ?", queryPostIds).
		Joins("Content").
		Joins("Author").
		Joins("MainCommunity").
		Find(&unorderPostList).Error

	if err != nil {
		return nil, errs.NewAppError(500, errs.CodeDBError, "数据库查询失败", err)
	}

	if len(unorderPostList) == 0 {
		return nil, errs.NewAppError(404, errs.CodePostListEmpty, "查询目标范围内帖子列表为空", nil)
	}

	// 按照 Redis 传进来的顺序组装结果
	for _, post := range unorderPostList {
		postIDMap[post.PostID] = post
	}

	for _, postID := range queryPostIds {
		if post, ok := postIDMap[postID]; ok {
			postList = append(postList, post)
		}
	}

	return postList, nil
}
