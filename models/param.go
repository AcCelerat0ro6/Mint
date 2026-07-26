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
