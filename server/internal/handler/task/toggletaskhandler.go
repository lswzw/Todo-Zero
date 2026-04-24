// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/internal/pkg/xerr"
	"server/internal/logic/task"
	"server/internal/svc"
	"server/internal/types"
)

func ToggleTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ToggleTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewToggleTaskLogic(r.Context(), svcCtx)
		resp, err := l.ToggleTask(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, xerr.SuccessResponse(resp))
		}
	}
}
