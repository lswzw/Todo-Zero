// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package stat

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/internal/pkg/xerr"
	"server/internal/logic/stat"
	"server/internal/svc"
)

func StatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := stat.NewStatLogic(r.Context(), svcCtx)
		resp, err := l.Stat()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, xerr.SuccessResponse(resp))
		}
	}
}
