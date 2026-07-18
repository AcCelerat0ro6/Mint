package models

// RegisterParam 注册请求参数
type RegisterParam struct {
	Username   string `json:"username"    binding:"required"`
	Password   string `json:"password"    binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email"       binding:"omitempty,email"`
	Gender     int8   `json:"gender"      binding:"omitempty,oneof=0 1 2"`
}
