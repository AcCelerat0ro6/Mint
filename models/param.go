package models

// 定义请求的参数结构体

// RegisterParam 注册请求参数
type RegisterParam struct {
	Username   string `json:"username"    binding:"required"`
	Password   string `json:"password"    binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email"       binding:"omitempty,email"`
	Gender     int8   `json:"gender"      binding:"omitempty,oneof=0 1 2"`
}

// LoginParam 登录请求参数
type LoginParam struct {
	Username string `json:"username"    binding:"required"`
	Password string `json:"password"    binding:"required"`
}

// RefreshTokenParam 刷新Token请求参数
type RefreshTokenParam struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreatePostParam 创建帖子请求参数
type CreatePostParam struct {
	CommunityID uint64 `json:"community_id" binding:"required"`
	Title       string `json:"title"        binding:"required"`
	Content     string `json:"content"      binding:"required"`
}

// GetPostListByTimeOrScoreParam 根据时间或分数获取帖子列表请求参数
type GetPostListByTimeOrScoreParam struct {
	Page        int    `form:"page"`
	Size        int    `form:"size"`
	Order       string `form:"order"        binding:"omitempty,oneof=time score"`
	CommunityID uint64 `form:"community_id"` // 新增：若为 0 查全局，不为 0 查指定社区
}
