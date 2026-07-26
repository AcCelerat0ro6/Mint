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

	api := r.Group("/api")

	v1 := api.Group("/v1")

	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome to mint blog !",
		})
	})

	// 注册路由
	v1.POST("/register", controller.RegisterHandler)

	// 登录路由
	v1.POST("/login", controller.LoginHandler)

	// 刷新Token路由
	v1.POST("/auth/refresh", controller.RefreshTokenHandler)

	// 获取社区列表路由
	v1.GET("/communitys", controller.GetCommunityListHandler)

	// 获取指定ID的社区路由
	v1.GET("/community/:id", controller.GetCommunityHandler)

	// 需要验证的路由组
	authGroup := v1.Group("/")
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

		// 创建帖子路由
		authGroup.POST("/createpost", controller.CreatePostHandler)
	}

	return r
}
