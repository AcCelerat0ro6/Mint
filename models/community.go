package models

import "time"

type CommunityInfo struct {
	CommunityID   uint64 `json:"community_id" gorm:"column:community_id;"`
	CommunityName string `json:"community_name" gorm:"column:community_name"`
}

func (CommunityInfo) TableName() string {
	return "community"
}

type CommunityDetail struct {
	CommunityID   uint64    `json:"community_id" gorm:"column:community_id;"`
	CommunityName string    `json:"community_name" gorm:"column:community_name"`
	Description   string    `json:"description" gorm:"column:description"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
}

func (CommunityDetail) TableName() string {
	return "community"
}
