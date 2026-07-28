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
	CodeSuccess            = 0
	CodeParamInvalid       = 40001 // 参数校验失败
	CodeUserExist          = 40002 // 用户名已存在
	CodeLoginUserNotExist  = 40003 // 登录用户不存在
	CodeLoginPasswordError = 40004 // 登录密码错误

	// JWT / 鉴权相关的错误码
	CodeAccessTokenExpired    = 40101 // Access Token 已过期
	CodeTokenSignatureInvalid = 40102 // Token 签名错误
	CodeTokenMalformed        = 40103 // Token 格式错误
	CodeTokenInvalid          = 40104 // Token 无效或未提供
	CodeRefreshTokenExpired   = 40105 // Refresh Token 已过期

	// 社区相关错误码
	CodeCommunityRecordNotFound = 40401 // 社区不存在

	// 帖子相关错误码
	CodePostRecordNotFound = 40601 // 帖子不存在
	CodePostListEmpty      = 40602 // 查询目标范围内帖子列表为空
	CodeLimitParamError    = 40603 // 分页参数格式错误

	// 解析整数错误码
	CodeParseUIntError = 40501 // 解析整数错误

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

	ErrLoginPasswordWrong = &AppError{
		HTTPCode: 400,
		BizCode:  CodeLoginPasswordError,
		Message:  "登录密码错误",
		Err:      nil,
	}

	// 预定义 JWT 相关错误
	ErrAccessTokenExpired = &AppError{
		HTTPCode: 401,
		BizCode:  CodeAccessTokenExpired,
		Message:  "Access Token 已过期，请调用刷新接口重新刷新",
		Err:      nil,
	}

	ErrRefreshTokenExpired = &AppError{
		HTTPCode: 401,
		BizCode:  CodeRefreshTokenExpired,
		Message:  "Refresh Token 已过期，请重新登录",
		Err:      nil,
	}

	ErrTokenSignatureInvalid = &AppError{
		HTTPCode: 401,
		BizCode:  CodeTokenSignatureInvalid,
		Message:  "Token 签名错误",
		Err:      nil,
	}

	ErrTokenMalformed = &AppError{
		HTTPCode: 401,
		BizCode:  CodeTokenMalformed,
		Message:  "Token 格式错误",
		Err:      nil,
	}

	ErrTokenInvalid = &AppError{
		HTTPCode: 401,
		BizCode:  CodeTokenInvalid,
		Message:  "Token 无效",
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
