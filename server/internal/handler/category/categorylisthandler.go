// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package category

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/internal/pkg/xerr"
	"server/internal/logic/category"
	"server/internal/svc"
)

func CategoryListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewCategoryListLogic(r.Context(), svcCtx)
		resp, err := l.CategoryList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, xerr.SuccessResponse(resp))
		}
	}
}
