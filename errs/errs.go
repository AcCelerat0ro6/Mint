package errs

import "fmt"

// AppError 统一的 API 错误结构
type AppError struct {
	HTTPCode int    // HTTP 状态码 (如 400, 500)
	BizCode  int    // 业务自定义码 (如 1004)
	Message  string // 给前端展示的友好提示
	Err      error  // 内部真实的底层错误 (仅给后端打日志用)
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("biz_code=%d, msg=%s, err=%v", e.BizCode, e.Message, e.Err)
	}
	return fmt.Sprintf("biz_code=%d, msg=%s", e.BizCode, e.Message)
}

// Unwrap 支持 errors.Is/As
func (e *AppError) Unwrap() error {
	return e.Err
}

// 业务错误码定义
const (
	CodeSuccess           = 0
	CodeParamInvalid      = 40001 // 参数校验失败
	CodeUserExist         = 40002 // 用户名已存在
	CodeDBError           = 50001 // 数据库查询错误
	CodeInternalError     = 50002 // 内部错误
	CodePasswordHashError = 50003 // 密码哈希错误
	CodeInsertDBError     = 50004 // 插入数据库错误
)

// 预定义的业务错误
var (
	ErrUserExist = &AppError{
		HTTPCode: 400,
		BizCode:  CodeUserExist,
		Message:  "用户名已存在",
		Err:      nil,
	}

	ErrDataBaseWrong = &AppError{
		HTTPCode: 500,
		BizCode:  CodeDBError,
		Message:  "数据库查询失败",
		Err:      nil,
	}
)

// NewAppError 创建新的 AppError
func NewAppError(httpCode, bizCode int, message string, err error) *AppError {
	return &AppError{
		HTTPCode: httpCode,
		BizCode:  bizCode,
		Message:  message,
		Err:      err,
	}
}
