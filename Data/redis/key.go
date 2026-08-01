package redis

// redis key 定义

const (
	KeyPrefix                       = "mint:"
	KeyPostTimeZSet                 = "post:time"             // zset;帖子及发帖时间
	KeyPostScoreZSet                = "post:score"            // zset;帖子及分数
	KeyPostVotedZSetPrefix          = "post:voted:"           // zset;帖子及投票用户
	KeyPostUpvoteCountHash          = "post:upvote:count"     // hash;帖子及点赞数量
	KeyCommunityPostTimeZSetPrefix  = "community:post:time:"  // zset;社区帖子发帖时间
	KeyCommunityPostScoreZSetPrefix = "community:post:score:" // zset;社区帖子分数
	KeyPostCommunityHash            = "post:community"        // hash;帖子ID映射社区ID
)
