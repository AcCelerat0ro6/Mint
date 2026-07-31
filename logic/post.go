package logic

import (
	"errors"
	"mint/Data/mysql"
	"mint/Data/redis"
	"mint/errs"
	"mint/models"
	"mint/pkg/snowflake"
	"strconv"
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

	// 4. 保存初始帖子分数和帖子时间到Redis
	if err := redis.SavePostScoreAndTime(strconv.FormatInt(postID, 10)); err != nil {
		if rollbackErr := mysql.DeletePostByID(uint64(postID)); rollbackErr != nil {
			return 0, errs.NewAppError(500, errs.CodeInternalError, "创建帖子失败，Redis 初始化失败且数据库补偿回滚失败", errors.Join(err, rollbackErr))
		}
		return 0, errs.NewAppError(500, errs.CodeInternalError, "创建帖子失败，Redis 初始化失败，数据库已回滚", err)
	}

	return uint64(postID), nil
}

func GetPostDetailByID(postID uint64) (*models.Post, error) {
	post, err := mysql.GetPostDetailByID(postID)
	if err != nil {
		return nil, err
	}

	if err := fillPostUpvoteCounts([]*models.Post{post}); err != nil {
		return nil, err
	}

	return post, nil
}

func GetPostList(page, size int) ([]models.Post, error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	postPointers := make([]*models.Post, 0, len(postList))
	for i := range postList {
		postPointers = append(postPointers, &postList[i])
	}

	if err := fillPostUpvoteCounts(postPointers); err != nil {
		return nil, err
	}

	return postList, nil
}

func GetPostListByTimeOrScore(param *models.GetPostListByTimeOrScoreParam) ([]*models.Post, error) {
	// 1. 从 Redis 获取帖子ID 列表
	queryPostIDs, err := redis.GetPostIDsInOrder(param)
	if err != nil {
		return nil, err
	}

	// 2. 从数据库查询帖子详情，返回的帖子按 Redis 传进来的顺序组装
	postList, err := mysql.GetPostListByIdList(queryPostIDs)
	if err != nil {
		return nil, err
	}

	if err := fillPostUpvoteCounts(postList); err != nil {
		return nil, err
	}

	return postList, nil
}

func fillPostUpvoteCounts(posts []*models.Post) error {
	if len(posts) == 0 {
		return errs.NewAppError(404, errs.CodePostListEmpty, "没有查询结果", nil)
	}

	postIDs := make([]string, 0, len(posts))
	for _, post := range posts {
		postIDs = append(postIDs, strconv.FormatUint(post.PostID, 10))
	}

	upvoteCounts, err := redis.GetPostUpvoteCounts(postIDs)
	if err != nil {
		return err
	}

	for i, post := range posts {
		post.UpvoteCount = upvoteCounts[i]
	}

	return nil
}
