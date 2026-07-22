package redis

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const refreshTokenKeyPrefix = "refresh_token:"

func getRefreshTokenKey(userID uint64) string {
	return fmt.Sprintf("%s%d", refreshTokenKeyPrefix, userID)
}

// SetRefreshToken 将 Refresh Token 的 jti 存入 Redis
func SetRefreshToken(userID uint64, jti string) error {
	expire := time.Duration(viper.GetInt("jwt.refresh_expire")) * time.Hour * 24
	if expire == 0 {
		expire = 7 * 24 * time.Hour
	}
	return rdb.Set(ctx, getRefreshTokenKey(userID), jti, expire).Err()
}

// CheckRefreshToken 检查 Refresh Token 是否存在且与 JTI 匹配
func CheckRefreshToken(userID uint64, jti string) (bool, error) {
	storedJTI, err := rdb.Get(ctx, getRefreshTokenKey(userID)).Result()
	if err != nil {
		return false, err // Redis 中不存在，可能是被踢下线或过期
	}
	if storedJTI == jti {
		return true, nil // 完全匹配
	}
	return false, nil // 不一致，说明在别处重新登录，当前 Token 被顶号作废
}

// DeleteRefreshToken 删除 Refresh Token（注销登录等使用）
func DeleteRefreshToken(userID uint64) error {
	return rdb.Del(ctx, getRefreshTokenKey(userID)).Err()
}
