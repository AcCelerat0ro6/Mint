package controller

import (
	"mint/logic"
	"mint/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostVoteController 帖子投票
func PostVoteController(c *gin.Context) {
	// 1. 参数校验
	var voteParam models.VoteDataParam
	if err := c.ShouldBindJSON(&voteParam); err != nil {
		c.Error(err)
		return
	}
	userID := c.GetUint64("userID")

	// 2. 投票逻辑
	err := logic.VoteForPost(strconv.FormatUint(userID, 10), voteParam.PostID, float64(*voteParam.Vote))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "vote success",
	})
}
