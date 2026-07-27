package controller

import (
	"mint/errs"
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

func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取帖子ID
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errs.NewAppError(400, errs.CodeParseUIntError, "帖子ID格式错误", err))
		return
	}

	// 2. 根据帖子ID查询帖子详情
	post, err := logic.GetPostDetailByID(postID)

	if err != nil {
		c.Error(err)
		return
	}

	// 3. 返回帖子详情
	c.JSON(http.StatusOK, gin.H{
		"message": "帖子详情查询成功",
		"post":    post,
	})
}
