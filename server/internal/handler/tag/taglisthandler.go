package tag

import (
	"net/http"

	"server/internal/logic/tag"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func TagListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tag.NewTagListLogic(r.Context(), svcCtx)
		resp, err := l.TagList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
