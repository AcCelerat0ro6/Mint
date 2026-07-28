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

func GetPostListHandler(c *gin.Context) {
	// 1. 获取参数并进行参数校验
	var (
		page, size int
		err        error
	)

	pagestr, sizestr := c.Query("page"), c.Query("size")

	if pagestr == "" {
		pagestr = "1"
	}

	if sizestr == "" {
		sizestr = "10"
	}

	if page, err = strconv.Atoi(pagestr); err != nil {
		c.Error(errs.NewAppError(400, errs.CodeLimitParamError, "分页参数格式错误", err))
		return
	}
	if page <= 0 {
		page = 1
	}

	if size, err = strconv.Atoi(sizestr); err != nil {
		c.Error(errs.NewAppError(400, errs.CodeLimitParamError, "分页参数格式错误", err))
		return
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	// 2. 获取帖子列表数据
	postList, err := logic.GetPostList(page, size)
	if err != nil {
		c.Error(err)
		return
	}

	// 2. 返回帖子列表
	c.JSON(http.StatusOK, gin.H{
		"message":  "帖子列表查询成功",
		"postList": postList,
	})
}
