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

func GetFocusWithIdByBodyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FocusRequestByBody
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := focus.NewGetFocusWithIdByBodyLogic(r.Context(), svcCtx)
		resp, err := l.GetFocusWithIdByBody(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
