package logic

import (
	"mint/Data/mysql"
	"mint/models"
	"mint/pkg/snowflake"
)

func CreatePost(post *models.CreatePostParam, userID uint64) (uint64, error) {
	// 1. 检查社区ID是否存在
	if err := mysql.FindCommunityExists(post.CommunityID); err != nil {
		return 0, err
	}

	// 2. 生成Post ID
	postID := snowflake.GenID()

	// 3. 保存帖子内容到数据库
	err := mysql.CreatePost(post, userID, uint64(postID))
	if err != nil {
		return 0, err
	}

	return uint64(postID), nil
}

func GetPostDetailByID(postID uint64) (*models.Post, error) {
	return mysql.GetPostDetailByID(postID)
}
