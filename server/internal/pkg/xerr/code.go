package xerr

import "fmt"

// 错误码
const (
	OK                 = 0
	ServerCommonError  = 10001
	RequestParamError  = 10002
	UserAlreadyExist   = 20001
	UserNotFoundError  = 20002
	PasswordError      = 20003
	UserDisabled       = 20004
	RegisterClosed     = 20005
	OldPasswordError   = 20006
	TaskNotFoundError     = 30001
	CategoryNotFoundError = 30002
	NoPermission         = 40001
	AdminRequired      = 40002
)

// 支持的语言
const (
	LangZhCN = "zh-CN"
	LangEn   = "en"
)

// 中文错误消息
var codeMessagesZhCN = map[int]string{
	OK:                "OK",
	ServerCommonError: "服务器内部错误",
	RequestParamError: "请求参数错误",
	UserAlreadyExist:  "用户名已存在",
	UserNotFoundError: "用户不存在",
	PasswordError:     "密码错误",
	UserDisabled:      "用户已被禁用",
	RegisterClosed:    "注册已关闭",
	OldPasswordError:  "原密码错误",
	TaskNotFoundError: "任务不存在",
	CategoryNotFoundError: "分类不存在",
	NoPermission:          "无权限操作",
	AdminRequired:     "需要管理员权限",
}

// 英文错误消息
var codeMessagesEn = map[int]string{
	OK:                "OK",
	ServerCommonError: "Internal server error",
	RequestParamError: "Request parameter error",
	UserAlreadyExist:  "Username already exists",
	UserNotFoundError: "User not found",
	PasswordError:     "Password error",
	UserDisabled:      "User has been disabled",
	RegisterClosed:    "Registration is closed",
	OldPasswordError:  "Old password error",
	TaskNotFoundError: "Task not found",
	CategoryNotFoundError: "Category not found",
	NoPermission:          "No permission",
	AdminRequired:     "Admin required",
}

// GetMessage 根据语言获取错误消息
func GetMessage(code int, lang string) string {
	var messages map[int]string
	switch lang {
	case LangEn:
		messages = codeMessagesEn
	default:
		messages = codeMessagesZhCN
	}
	msg, ok := messages[code]
	if !ok {
		if lang == LangEn {
			return "Unknown error"
		}
		return "未知错误"
	}
	return msg
}

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// NewCodeError 创建错误（默认中文）
func NewCodeError(code int) *CodeError {
	return NewCodeErrorWithLang(code, LangZhCN)
}

// NewCodeErrorWithLang 创建指定语言的错误
func NewCodeErrorWithLang(code int, lang string) *CodeError {
	return &CodeError{Code: code, Msg: GetMessage(code, lang)}
}

func NewCodeErrFromMsg(msg string) *CodeError {
	return &CodeError{Code: ServerCommonError, Msg: msg}
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

// ValidationError 表示输入验证失败的安全错误，其消息可以安全地返回给客户端。
// 与 CodeError 不同，ValidationError 携带的是动态的、上下文相关的验证提示，
// 如 "分类名称最多20个字符"。所有非 ValidationError 的普通 error 均不会暴露原始消息。
type ValidationError struct {
	Msg string
}

// NewValidationError 创建一个可安全暴露给客户端的验证错误。
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Msg: msg}
}

func (e *ValidationError) Error() string {
	return e.Msg
}

// RateLimitError 表示请求频率超过限制的错误，返回 HTTP 429 状态码。
type RateLimitError struct {
	Msg string
}

// NewRateLimitError 创建一个限流错误。
func NewRateLimitError(msg string) *RateLimitError {
	return &RateLimitError{Msg: msg}
}

func (e *RateLimitError) Error() string {
	return e.Msg
}
