package routers

import (
	"mint/controller"
	"mint/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), controller.ErrorMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome to mint blog !",
		})
	})

	// 注册路由
	r.POST("/register", controller.RegisterHandler)

	return r
}
