package redis

import (
	"errors"
	"mint/errs"
	"mint/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func SavePostScoreAndTime(postID string) error {
	curtime := time.Now()

	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// 这里的 ZAdd 不会立即发往 Redis，而是缓存在管道中
		pipe.ZAdd(ctx, KeyPrefix+KeyPostTimeZSet, redis.Z{
			Score:  float64(curtime.Unix()),
			Member: postID,
		})

		pipe.ZAdd(ctx, KeyPrefix+KeyPostScoreZSet, redis.Z{
			Score:  float64(curtime.Unix()),
			Member: postID,
		})

		pipe.HSet(ctx, KeyPrefix+KeyPostUpvoteCountHash, postID, 0)
		return nil // 返回 nil 表示管道组装成功
	})
	if err != nil {
		return errs.NewAppError(500, errs.CodeRedisZSetInsertError, "插入帖子分数到 ZSet 失败", err)
	}
	return nil
}

func GetPostIDsInOrder(param *models.GetPostListByTimeOrScoreParam) ([]string, error) {
	var queryKey string
	if param.Order == "time" {
		queryKey = KeyPrefix + KeyPostTimeZSet
	} else {
		queryKey = KeyPrefix + KeyPostScoreZSet
	}

	startIndex := (param.Page - 1) * param.Size
	endIndex := startIndex + param.Size - 1

	postIDs, err := rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   queryKey,
		Start: startIndex,
		Stop:  endIndex,
		Rev:   true,
	}).Result()

	if err != nil {
		return nil, errs.NewAppError(500, errs.CodeRedisQueryError, "Redis查询 ID 列表失败", err)
	}
	if len(postIDs) == 0 {
		return nil, errs.NewAppError(404, errs.CodePostListEmpty, "查询目标范围内帖子列表为空", nil)
	}

	return postIDs, nil
}

func GetPostUpvoteCounts(postIDs []string) ([]int64, error) {
	upvoteCounts := make([]int64, len(postIDs))
	if len(postIDs) == 0 {
		return upvoteCounts, nil
	}

	upvoteSlice, err := rdb.HMGet(ctx, KeyPrefix+KeyPostUpvoteCountHash, postIDs...).Result()
	if err != nil {
		return nil, errs.NewAppError(500, errs.CodeRedisQueryError, "Redis查询帖子点赞数失败", err)
	}

	for i := range postIDs {
		if upvoteSlice[i] == nil {
			// 兼容历史数据或异常回滚场景，缺失的点赞数按 0 处理
			upvoteCounts[i] = 0
			continue
		}

		upvoteCountStr, ok := upvoteSlice[i].(string)
		if !ok {
			return nil, errs.NewAppError(500, errs.CodeRedisQueryError, "Redis点赞数类型错误", errors.New("HMGet 返回的点赞数不是 string 类型"))
		}

		upvoteCounts[i], err = strconv.ParseInt(upvoteCountStr, 10, 64)
		if err != nil {
			return nil, errs.NewAppError(500, errs.CodeRedisQueryError, "解析帖子点赞数失败", err)
		}
	}

	return upvoteCounts, nil
}
