package xerr

import (
	"context"
	"log"
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
		msg := GetMessage(e.Code, lang)
		return http.StatusOK, &Body{Code: e.Code, Msg: msg}
	default:
		msg := err.Error()
		if msg == "登录尝试次数过多，请稍后再试" || msg == "Too many login attempts, please try again later" {
			if lang == LangEn {
				return http.StatusTooManyRequests, &Body{Code: 42901, Msg: "Too many login attempts, please try again later"}
			}
			return http.StatusTooManyRequests, &Body{Code: 42901, Msg: msg}
		}
		if msg != "" {
			return http.StatusOK, &Body{Code: RequestParamError, Msg: msg}
		}
		log.Printf("[xerr] Internal error: %v", err)
		return http.StatusInternalServerError, &Body{Code: ServerCommonError, Msg: GetMessage(ServerCommonError, lang)}
	}
}
