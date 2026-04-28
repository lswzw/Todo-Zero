package tag

import (
	"net/http"

	"server/internal/logic/tag"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteTagReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tag.NewDeleteTagLogic(r.Context(), svcCtx)
		resp, err := l.DeleteTag(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
