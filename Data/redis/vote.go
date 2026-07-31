package redis

import (
	_ "embed"
	"errors"
	"mint/errs"

	"github.com/redis/go-redis/v9"
)

func GetPostCreateTime(postID string) (float64, error) {
	postTime, err := rdb.ZScore(ctx, KeyPrefix+KeyPostTimeZSet, postID).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, errs.NewAppError(404, errs.CodeRedisZSetMemberNotExist, "Redis中找不到帖子相关的记录", err)
		} else {
			return 0, errs.NewAppError(500, errs.CodeRedisQueryError, "Redis查询错误", err)
		}
	}

	return postTime, nil
}

//go:embed vote.lua
var voteLuaScript string

var voteScript = redis.NewScript(voteLuaScript)

func VoteForPost(userID string, postID string, curValue float64) error {
	personalKey := KeyPrefix + KeyPostVotedZSetPrefix + postID
	scoreKey := KeyPrefix + KeyPostScoreZSet
	upvoteCountKey := KeyPrefix + KeyPostUpvoteCountHash

	keys := []string{personalKey, scoreKey, upvoteCountKey}
	args := []interface{}{userID, curValue, postID}

	// 运行 Lua 脚本，保证原子性
	_, err := voteScript.Run(ctx, rdb, keys, args...).Result()
	if err != nil {
		return errs.NewAppError(500, errs.CodeRedisModifyZSetMemberError, "处理点赞并发逻辑失败", err)
	}

	return nil
}
