// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package focus

import (
	"gozeroapi/internal/logic/focus"
	"net/http"

	"gozeroapi/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFocusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := focus.NewGetFocusLogic(r.Context(), svcCtx)
		resp, err := l.GetFocus()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
