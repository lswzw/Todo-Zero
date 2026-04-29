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
	case *ValidationError:
		// 验证错误：消息已经过审查，可安全暴露给客户端
		return http.StatusOK, &Body{Code: RequestParamError, Msg: e.Msg}
	case *RateLimitError:
		// 限流错误：消息可安全暴露，返回 HTTP 429
		return http.StatusTooManyRequests, &Body{Code: 42901, Msg: e.Msg}
	default:
		// 未知错误：不暴露内部错误详情，统一返回通用消息
		log.Printf("[xerr] Internal error: %v", err)
		return http.StatusInternalServerError, &Body{Code: ServerCommonError, Msg: GetMessage(ServerCommonError, lang)}
	}
}
