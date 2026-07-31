package models

import "time"

type Post struct {
	PostID        uint64           `gorm:"column:post_id;primaryKey"                     json:"post_id,string"`
	AuthorID      uint64           `gorm:"column:author_id"                              json:"-"`
	CommunityID   uint64           `gorm:"column:community_id"                           json:"-"`
	Status        uint8            `gorm:"column:status"                                 json:"status"`
	Title         string           `gorm:"column:title"                                  json:"title"`
	UpvoteCount   int64            `gorm:"-"                                             json:"upvote_count"`
	CreatedAt     time.Time        `gorm:"column:created_at"                             json:"created_at"`
	UpdatedAt     time.Time        `gorm:"column:updated_at"                             json:"updated_at"`
	Content       *PostContent     `gorm:"foreignKey:PostID;references:PostID"           json:"content"`
	Author        *Author          `gorm:"foreignKey:AuthorID;references:UserID"         json:"author"`
	MainCommunity *CommunityDetail `gorm:"foreignKey:CommunityID;references:CommunityID" json:"main_community"`
}

func (Post) TableName() string {
	return "post"
}

type PostContent struct {
	PostID  uint64 `gorm:"column:post_id;primaryKey"     json:"-"`
	Content string `gorm:"column:content"                json:"content"`
}

func (PostContent) TableName() string {
	return "post_content"
}
