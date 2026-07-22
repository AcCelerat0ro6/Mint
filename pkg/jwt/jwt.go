package jwt

import (
	"errors"
	"mint/errs"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	mySecret = []byte("备用密钥:夏天夏天悄悄过去")
)

// AccessClaims 仅保留 Access Token 必需的标准字段。
type AccessClaims struct {
	jwt.RegisteredClaims
}

// UserID 从 sub 中解析当前 token 对应的 userID。
func (c *AccessClaims) UserID() (uint64, error) {
	return parseUserID(c.Subject)
}

// RefreshClaims 仅保留 Refresh Token 必需的标准字段。
type RefreshClaims struct {
	jwt.RegisteredClaims
}

// UserID 从 sub 中解析当前 token 对应的 userID。
func (c *RefreshClaims) UserID() (uint64, error) {
	return parseUserID(c.Subject)
}

// Init 初始化JWT，从配置读取密钥等信息（可选）
func Init() {
	if secret := viper.GetString("jwt.secret"); secret != "" {
		mySecret = []byte(secret)
	}
}

// GenTokens 生成 Access Token 和 Refresh Token
func GenTokens(userID uint64) (aToken, rToken, jti string, err error) {
	// 1. 生成 Access Token (有效期短，如 15-30 分钟)
	aExp := time.Now().Add(time.Duration(viper.GetInt("jwt.access_expire")) * time.Minute)
	if viper.GetInt("jwt.access_expire") == 0 {
		aExp = time.Now().Add(30 * time.Minute)
	}
	c := AccessClaims{
		jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(userID, 10),
			ExpiresAt: jwt.NewNumericDate(aExp),
			Issuer:    "mint",
		},
	}
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	if err != nil {
		return "", "", "", err
	}

	// 2. 生成 Refresh Token (有效期长，如 7 天)，带有 JTI
	rExp := time.Now().Add(time.Duration(viper.GetInt("jwt.refresh_expire")) * time.Hour * 24)
	if viper.GetInt("jwt.refresh_expire") == 0 {
		rExp = time.Now().Add(7 * 24 * time.Hour)
	}
	jti = uuid.New().String()
	rc := RefreshClaims{
		jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(userID, 10),
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(rExp),
			Issuer:    "mint",
		},
	}
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rc).SignedString(mySecret)
	return
}

// ParseToken 解析 Access Token
func ParseToken(tokenString string) (*AccessClaims, error) {
	var mc = new(AccessClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.ErrAccessTokenExpired
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errs.ErrTokenSignatureInvalid
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errs.ErrTokenMalformed
		}
		return nil, errs.NewAppError(401, errs.CodeTokenInvalid, "Token 解析失败", err)
	}
	if token.Valid {
		return mc, nil
	}
	return nil, errs.ErrTokenInvalid
}

// ParseRefreshToken 解析 Refresh Token
func ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	var rc = new(RefreshClaims)
	token, err := jwt.ParseWithClaims(tokenString, rc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.ErrRefreshTokenExpired
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errs.ErrTokenSignatureInvalid
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errs.ErrTokenMalformed
		}
		return nil, errs.NewAppError(401, errs.CodeTokenInvalid, "Token 解析失败", err)
	}
	if token.Valid {
		return rc, nil
	}
	return nil, errs.ErrTokenInvalid
}

func parseUserID(subject string) (uint64, error) {
	userID, err := strconv.ParseUint(subject, 10, 64)
	if err != nil {
		return 0, errs.NewAppError(401, errs.CodeTokenInvalid, "Token 中的 sub 非法", err)
	}
	return userID, nil
}
