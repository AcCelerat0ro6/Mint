package controller

import (
	"mint/logic"
	"mint/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	// 1. 参数校验
	param := new(models.RegisterParam)
	if err := c.ShouldBindJSON(param); err != nil {
		c.Error(err)
		return
	}
	// 2. 业务处理
	if err := logic.SignUp(param); err != nil {
		c.Error(err)
		return
	}

	// 3. 返回响应
	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
	})
}
