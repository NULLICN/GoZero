// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"client/m4api/internal/logic"
	"client/m4api/internal/svc"
	"client/m4api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SayHelloPostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SayHelloReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSayHelloPostLogic(r.Context(), svcCtx)
		resp, err := l.SayHelloPost(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
