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

// getLangFromCtx 从 context 中获取语言（通过 locale middleware 设置）
func getLangFromCtx(ctx context.Context) string {
	lang, ok := ctx.Value(LocaleKey).(string)
	if !ok {
		return LangZhCN
	}
	return lang
}

// ErrorResponse 错误响应
func ErrorResponse(ctx context.Context, err error) (int, interface{}) {
	lang := getLangFromCtx(ctx)

	switch e := err.(type) {
	case *CodeError:
		// 根据语言重新获取错误消息
		msg := GetMessage(e.Code, lang)
		return http.StatusOK, &Body{Code: e.Code, Msg: msg}
	default:
		msg := err.Error()
		// 登录限流错误
		if msg == "登录尝试次数过多，请稍后再试" {
			if lang == LangEn {
				return http.StatusTooManyRequests, &Body{Code: 42901, Msg: "Too many login attempts, please try again later"}
			}
			return http.StatusTooManyRequests, &Body{Code: 42901, Msg: msg}
		}
		// Validate() 等返回的普通 error 视为参数错误
		if msg != "" {
			return http.StatusOK, &Body{Code: RequestParamError, Msg: msg}
		}
		return http.StatusInternalServerError, &Body{Code: ServerCommonError, Msg: GetMessage(ServerCommonError, lang)}
	}
}
