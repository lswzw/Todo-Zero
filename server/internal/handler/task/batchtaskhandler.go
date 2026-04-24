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

func BatchTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewBatchTaskLogic(r.Context(), svcCtx)
		resp, err := l.BatchTask(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, xerr.SuccessResponse(resp))
		}
	}
}
