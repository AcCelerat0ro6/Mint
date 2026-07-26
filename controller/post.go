package controller

import (
	"mint/logic"
	"mint/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	// 1. 获取帖子参数及参数校验
	var post models.CreatePostParam
	if err := c.ShouldBindJSON(&post); err != nil {
		c.Error(err)
		return
	}
	// 调用逻辑层创建帖子
	var postID uint64
	var err error
	if postID, err = logic.CreatePost(&post, c.GetUint64("userID")); err != nil {
		c.Error(err)
		return
	}

	// 返回创建成功的响应
	c.JSON(http.StatusOK, gin.H{
		"message": "帖子创建成功",
		"postID":  strconv.FormatUint(postID, 10), //为了防止精度损失，转换成字符串再传递给前端,
	})
}
