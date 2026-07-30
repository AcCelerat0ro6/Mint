package redis

import (
	"mint/errs"
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
		return nil // 返回 nil 表示管道组装成功
	})
	if err != nil {
		return errs.NewAppError(500, errs.CodeRedisZSetInsertError, "插入帖子分数到 ZSet 失败", err)
	}
	return nil
}
