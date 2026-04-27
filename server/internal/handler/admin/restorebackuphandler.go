package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/internal/logic/admin"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"
)

func RestoreBackupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RestoreBackupReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.NewCodeError(xerr.RequestParamError))
			return
		}

		l := admin.NewRestoreBackupLogic(r.Context(), svcCtx)
		resp, err := l.RestoreBackup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, xerr.SuccessResponse(resp))
		}
	}
}
