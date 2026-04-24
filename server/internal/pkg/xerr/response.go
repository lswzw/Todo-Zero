package xerr

import (
	"context"
	"net/http"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) *Body {
	return &Body{Code: OK, Msg: "ok", Data: data}
}

// ErrorResponse 错误响应
func ErrorResponse(ctx context.Context, err error) (int, interface{}) {
	switch e := err.(type) {
	case *CodeError:
		return http.StatusOK, &Body{Code: e.Code, Msg: e.Msg}
	default:
		msg := err.Error()
		// 登录限流错误
		if msg == "登录尝试次数过多，请稍后再试" {
			return http.StatusTooManyRequests, &Body{Code: 42901, Msg: msg}
		}
		// Validate() 等返回的普通 error 视为参数错误
		if msg != "" {
			return http.StatusOK, &Body{Code: RequestParamError, Msg: msg}
		}
		return http.StatusInternalServerError, &Body{Code: ServerCommonError, Msg: "服务器内部错误"}
	}
}
