package middleware

import (
	"net/http"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (m *AdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin, err := jwtx.GetIsAdminFromCtx(r.Context())
		if err != nil || isAdmin != 1 {
			logx.Errorf("[AdminMiddleware] non-admin access attempt: %s %s", r.Method, r.URL.Path)
			httpx.ErrorCtx(r.Context(), w, xerr.NewCodeError(xerr.AdminRequired))
			return
		}
		next(w, r)
	}
}
