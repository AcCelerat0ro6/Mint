package controller

import (
	"mint/errs"
	"mint/logic"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCommunityListHandler获取社区列表
func GetCommunityListHandler(c *gin.Context) {
	// 从数据库中查询社区列表
	communityList, err := logic.GetCommunityList()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "community list",
		"data":    communityList,
	})
}

// GetCommunityHandler获取指定ID的社区
func GetCommunityHandler(c *gin.Context) {
	// 从URL参数中获取社区ID
	communityIDstr := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDstr, 10, 64)
	if err != nil {
		c.Error(&errs.AppError{
			HTTPCode: 400,
			BizCode:  errs.CodeParseUIntError,
			Message:  "传入的ID不是一个UNSIGNED INT类型",
			Err:      err,
		})
		return
	}

	// 调用逻辑层获取社区详情
	community, err := logic.GetCommunityByID(communityID)
	if err != nil {
		c.Error(err)
		return
	}

	// 返回社区详情
	c.JSON(http.StatusOK, gin.H{
		"message":   "community detail",
		"community": community,
	})
}
