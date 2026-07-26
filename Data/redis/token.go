package redis

import (
	"errors"
	"fmt"
	"mint/errs"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const refreshTokenKeyPrefix = "refresh_token:"

type RefreshTokenStatus int

const (
	RefreshTokenStatusUnknown RefreshTokenStatus = iota
	RefreshTokenStatusValid
	RefreshTokenStatusExpired
	RefreshTokenStatusInvalidated
)

func getRefreshTokenKey(userID uint64) string {
	return fmt.Sprintf("%s%d", refreshTokenKeyPrefix, userID)
}

// SetRefreshToken 将 Refresh Token 的 jti 存入 Redis
func SetRefreshToken(userID uint64, jti string) error {
	expire := time.Duration(viper.GetInt("jwt.refresh_expire")) * time.Hour * 24
	if expire == 0 {
		expire = 7 * 24 * time.Hour
	}
	if err := rdb.Set(ctx, getRefreshTokenKey(userID), jti, expire).Err(); err != nil {
		return errs.NewAppError(500, errs.CodeInternalError, "保存 Refresh Token 登录态失败", err)
	}
	return nil
}

// CheckRefreshToken 检查 Refresh Token 是否存在且与 JTI 匹配
func CheckRefreshToken(userID uint64, jti string) (RefreshTokenStatus, error) {
	storedJTI, err := rdb.Get(ctx, getRefreshTokenKey(userID)).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return RefreshTokenStatusExpired, nil
		}
		return RefreshTokenStatusUnknown, errs.NewAppError(500, errs.CodeInternalError, "读取 Refresh Token 登录态失败", err)
	}
	if storedJTI == jti {
		return RefreshTokenStatusValid, nil // 完全匹配
	}
	return RefreshTokenStatusInvalidated, nil // 不一致，说明在别处重新登录，当前 Token 被顶号作废
}

// DeleteRefreshToken 删除 Refresh Token（注销登录等使用）
func DeleteRefreshToken(userID uint64) error {
	if err := rdb.Del(ctx, getRefreshTokenKey(userID)).Err(); err != nil {
		return errs.NewAppError(500, errs.CodeInternalError, "删除 Refresh Token 登录态失败", err)
	}
	return nil
}
