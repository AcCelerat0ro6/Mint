package logic

import (
	"mint/Data/redis"
	"mint/errs"
	"time"
)

// 投票功能:
// 1. 用户投票的数据
// 2. 取消投票
// 3. 查看投票结果

func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 去redis取帖子发布时间
	postTime, err := redis.GetPostCreateTime(postID)
	if err != nil {
		return err
	}

	// 判断时间是否超过一周，如果超过一周，我们不再给予其票数上的改变
	if float64(time.Now().Unix())-postTime > 7*24*60*60 {
		return errs.NewAppError(400, errs.CodeVoteExpiredPost, "帖子发布时间超过一周，不能投票", nil)
	}

	// 2. 调用 Redis Lua 脚本原子化处理投票并发逻辑 (查历史记录、更新总分、更新个人记录)
	if err := redis.VoteForPost(userID, postID, value); err != nil {
		return err
	}

	return nil
}
