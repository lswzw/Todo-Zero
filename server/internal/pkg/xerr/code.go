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

// 错误消息
var codeMessages = map[int]string{
	OK:                "OK",
	ServerCommonError: "服务器内部错误",
	RequestParamError: "请求参数错误",
	UserAlreadyExist:  "用户名已存在",
	UserNotFoundError: "用户不存在",
	PasswordError:     "密码错误",
	UserDisabled:      "用户已被禁用",
	RegisterClosed:    "注册已关闭",
	OldPasswordError:  "原密码错误",
	TaskNotFoundError:     "任务不存在",
	CategoryNotFoundError: "分类不存在",
	NoPermission:          "无权限操作",
	AdminRequired:     "需要管理员权限",
}

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewCodeError(code int) *CodeError {
	msg, ok := codeMessages[code]
	if !ok {
		msg = "未知错误"
	}
	return &CodeError{Code: code, Msg: msg}
}

func NewCodeErrFromMsg(msg string) *CodeError {
	return &CodeError{Code: ServerCommonError, Msg: msg}
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}
