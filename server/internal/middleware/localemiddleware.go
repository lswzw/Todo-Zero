package middleware

import (
	"context"
	"net/http"
	"strings"

	"server/internal/pkg/xerr"
)

// LocaleMiddleware 语言中间件，解析 Accept-Language 并存储到 context 中
type LocaleMiddleware struct{}

func NewLocaleMiddleware() *LocaleMiddleware {
	return &LocaleMiddleware{}
}

func (m *LocaleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := parseAcceptLanguage(r.Header.Get("Accept-Language"))
		ctx := context.WithValue(r.Context(), xerr.LocaleKey, lang)
		next(w, r.WithContext(ctx))
	}
}

// parseAcceptLanguage 解析 Accept-Language 头
func parseAcceptLanguage(acceptLang string) string {
	if acceptLang == "" {
		return xerr.LangZhCN
	}

	parts := strings.Split(acceptLang, ",")
	for _, part := range parts {
		lang := strings.TrimSpace(strings.Split(part, ";")[0])
		if strings.HasPrefix(lang, "en") {
			return xerr.LangEn
		} else if strings.HasPrefix(lang, "zh") {
			return xerr.LangZhCN
		}
	}

	return xerr.LangZhCN
}

// GetLangFromCtx 从 context 中获取语言
func GetLangFromCtx(ctx context.Context) string {
	lang, ok := ctx.Value(xerr.LocaleKey).(string)
	if !ok {
		return xerr.LangZhCN
	}
	return lang
}
