// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package focus

import (
	"gozeroapi/internal/logic/focus"
	"net/http"

	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFocusWithIdByQueryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FocusRequestByQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := focus.NewGetFocusWithIdByQueryLogic(r.Context(), svcCtx)
		resp, err := l.GetFocusWithIdByQuery(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
