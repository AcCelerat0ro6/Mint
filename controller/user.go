package controller

import (
	"mint/Data/redis"
	"mint/errs"
	"mint/logic"
	"mint/models"
	"mint/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandler用户注册函数
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

// LoginHandler用户登录函数
func LoginHandler(c *gin.Context) {
	// 1. 参数校验
	param := new(models.LoginParam)
	if err := c.ShouldBindJSON(param); err != nil {
		c.Error(err)
		return
	}
	// 2. 用户登录，获取 tokens(access token 和 refresh token)
	aToken, rToken, err := logic.Login(param)
	if err != nil {
		c.Error(err)
		return
	}

	// 3. 返回响应
	c.JSON(http.StatusOK, gin.H{
		"message":       "login success",
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}

// RefreshTokenHandler 刷新 Access Token
func RefreshTokenHandler(c *gin.Context) {
	// 1. 获取参数
	param := new(models.RefreshTokenParam)
	if err := c.ShouldBindJSON(param); err != nil {
		c.Error(err)
		return
	}

	// 2. 解析 Refresh Token
	rc, err := jwt.ParseRefreshToken(param.RefreshToken)
	if err != nil {
		c.Error(err) // 解析失败直接抛出 AppError 给全局处理
		return
	}

	userID, err := rc.UserID()
	if err != nil {
		c.Error(err)
		return
	}
	jti := rc.ID

	// 3. 检查 Redis 中该 Refresh Token 的 JTI 是否匹配
	status, err := redis.CheckRefreshToken(userID, jti)
	if err != nil {
		c.Error(err)
		return
	}
	if status == redis.RefreshTokenStatusExpired {
		c.Error(errs.ErrRefreshTokenExpired)
		return
	}
	if status == redis.RefreshTokenStatusInvalidated {
		c.Error(errs.NewAppError(http.StatusUnauthorized, errs.CodeTokenInvalid, "Refresh Token 已失效、被拉黑或在别处重新登录", nil))
		return
	}

	// 4. 生成新的 Access Token 和 Refresh Token (Token轮转)
	newAToken, newRToken, newJTI, err := jwt.GenTokens(userID)
	if err != nil {
		c.Error(errs.NewAppError(http.StatusInternalServerError, errs.CodeInternalError, "生成新Token失败", err))
		return
	}

	// 5. 将新的 JTI 存入 Redis，实现 Token 轮转
	if err := redis.SetRefreshToken(userID, newJTI); err != nil {
		c.Error(err)
		return
	}

	// 6. 返回新的 Token
	c.JSON(http.StatusOK, gin.H{
		"message":       "refresh success",
		"access_token":  newAToken,
		"refresh_token": newRToken,
	})
}
