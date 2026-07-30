package redis

// redis key 定义

const (
	KeyPrefix              = "mint:"
	KeyPostTimeZSet        = "post:time"   // zset;帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  // zset;帖子及分数
	KeyPostVotedZSetPrefix = "post:voted:" // zset;帖子及投票用户
)
