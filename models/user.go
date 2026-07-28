package models

import "time"

// 定义表模型

// User 用户模型
type User struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id,string"`
	UserID    uint64    `gorm:"column:user_id" json:"user_id,string"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"-"`
	Email     *string   `gorm:"column:email" json:"email"`
	Gender    int8      `gorm:"column:gender" json:"gender"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (u *User) TableName() string {
	return "user"
}

// 登录用户表模型
type LoginUser struct {
	ID       uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	UserID   uint64 `gorm:"column:user_id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

// 社区作者用户表模型
type Author struct {
	UserID   uint64  `gorm:"column:user_id"     json:"author_id,string"`
	Username string  `gorm:"column:username"    json:"author_name"`
	Email    *string `gorm:"column:email"       json:"author_email"`
	Gender   int8    `gorm:"column:gender"      json:"author_gender"`
}

func (a *Author) TableName() string {
	return "user"
}
