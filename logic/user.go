package logic

import (
	"fmt"
	"mint/Data/mysql"
	"mint/errs"
	"mint/models"
	"mint/pkg/snowflake"

	"golang.org/x/crypto/bcrypt"
)

// SignUp判断用户是否存在，不存在则注册
func SignUp(param *models.RegisterParam) error {
	// 1. 判断用户是否存在
	exist, err := mysql.CheckUserExist(param.Username)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrDataBaseWrong, err)
	}
	if exist {
		return errs.ErrUserExist
	}

	// 2. 生成UID
	uid := snowflake.GenID()

	// 3. 密码加密
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		return errs.NewAppError(500, errs.CodePasswordHashError, "密码哈希失败", err)
	}

	// 4. 保存用户信息到数据库
	// 4.1 构造用户实例
	var email *string
	if param.Email != "" {
		email = &param.Email
	}
	user := &models.User{
		UserID:   uint64(uid),
		Username: param.Username,
		Password: string(hashPassword),
		Email:    email,
		Gender:   param.Gender,
	}

	// 4.2 保存用户到数据库
	return mysql.InsertUser(user)
}
