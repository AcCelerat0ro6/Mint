package models

import "time"

type Post struct {
	PostID      uint64    `gorm:"column:post_id;primaryKey"`
	AuthorID    uint64    `gorm:"column:author_id"`
	CommunityID uint64    `gorm:"column:community_id"`
	Status      uint8     `gorm:"column:status"`
	Title       string    `gorm:"column:title"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Post) TableName() string {
	return "post"
}

type PostContent struct {
	PostID  uint64 `gorm:"column:post_id;primaryKey"`
	Content string `gorm:"column:content"`
}

func (PostContent) TableName() string {
	return "post_content"
}
