package controller

import (
	"errors"
	"mint/errs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"mint/pkg/jwt"
)

// ErrorMiddleware 全局错误处理中间件
// 将错误分为两类：
// 1. 客户端错误（400）：参数校验错误 + 业务错误（如用户名已存在）
// 2. 系统错误（500）：数据库错误等内部错误，记录日志
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var (
			validationMsgs []string
			businessErr    *errs.AppError
			hasSystemErr   bool
		)

		// 遍历所有错误，分为两类
		for _, ginErr := range c.Errors {
			err := ginErr.Err

			// 1. 参数校验错误 -> 归为客户端错误（400）
			var valErrs validator.ValidationErrors
			if errors.As(err, &valErrs) {
				translated := valErrs.Translate(Trans)
				for _, msg := range translated {
					validationMsgs = append(validationMsgs, msg)
				}
				continue
			}

			// 2. AppError
			var appErr *errs.AppError
			if errors.As(err, &appErr) {
				if appErr.BizCode >= 50000 {
					// 系统错误（BizCode >= 50000）：记录日志，标记为 500
					hasSystemErr = true
					zap.L().Error("系统错误",
						zap.Int("biz_code", appErr.BizCode),
						zap.String("message", appErr.Message),
						zap.Error(appErr.Err),
						zap.String("path", c.Request.URL.Path),
					)
				} else {
					// 业务错误（BizCode < 50000）：归为客户端错误（400），保存第一个
					if businessErr == nil {
						businessErr = appErr
					}
				}
				continue
			}

			// 3. 其他未知错误 -> 归为系统错误（500）
			hasSystemErr = true
			zap.L().Error("未知系统错误",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
			)
		}

		// 返回逻辑：系统错误优先（500），否则业务错误和参数校验错误都返回 400
		if hasSystemErr {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": errs.CodeInternalError,
				"msg":  "服务器内部异常，请稍后再试",
			})
			return
		}

		// 业务错误和参数校验错误都返回 400
		if businessErr != nil {
			c.JSON(businessErr.HTTPCode, gin.H{
				"code": businessErr.BizCode,
				"msg":  businessErr.Message,
			})
			return
		}

		if len(validationMsgs) > 0 {
			errMsg := strings.Join(validationMsgs, "; ")
			c.JSON(http.StatusBadRequest, gin.H{
				"code": errs.CodeParamInvalid,
				"msg":  "参数校验失败",
				"data": errMsg,
			})
			return
		}
	}
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Error(errs.NewAppError(http.StatusUnauthorized, errs.CodeTokenInvalid, "请求头中auth为空", nil))
			c.Abort()
			return
		}

		// 2. 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Error(errs.NewAppError(http.StatusUnauthorized, errs.CodeTokenInvalid, "请求头中auth格式有误", nil))
			c.Abort()
			return
		}

		// 3. 解析 Token
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			// 直接将 pkg/jwt 返回的 AppError 送入全局错误处理中间件
			c.Error(err)
			c.Abort()
			return
		}

		userID, err := mc.UserID()
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// 4. 将当前请求的 userID 信息保存到请求的上下文 c 上
		c.Set("userID", userID)
		c.Next()
	}
}
