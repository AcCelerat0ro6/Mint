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

	// 登录路由
	r.POST("/login", controller.LoginHandler)

	// 刷新Token路由
	r.POST("/api/auth/refresh", controller.RefreshTokenHandler)

	// 未验证的pong路由
	r.GET("/pong/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong public"})
	})

	// 需要验证的路由组
	authGroup := r.Group("/")
	authGroup.Use(controller.JWTAuthMiddleware())
	{
		// 验证的pong路由
		authGroup.GET("/pong", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{
				"message": "pong authenticated",
				"userID":  userID,
			})
		})
	}

	return r
}
