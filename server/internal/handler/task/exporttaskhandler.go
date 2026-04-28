package task

import (
	"net/http"

	"server/internal/logic/task"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ExportTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := task.NewExportTaskLogic(r.Context(), svcCtx)
		err := l.ExportTask(&req, w)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
	}
}
