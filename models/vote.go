package models

type VoteDataParam struct {
	PostID string `json:"postID"   binding:"required"`
	Vote   int    `json:"vote"            binding:"required,oneof=1 -1 0"`
}
