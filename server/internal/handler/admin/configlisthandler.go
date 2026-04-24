// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/internal/pkg/xerr"
	"server/internal/logic/admin"
	"server/internal/svc"
)

func ConfigListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewConfigListLogic(r.Context(), svcCtx)
		resp, err := l.ConfigList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, xerr.SuccessResponse(resp))
		}
	}
}
