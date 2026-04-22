// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"firstdemo/internal/logic"
	"firstdemo/internal/svc"
	"firstdemo/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FirstdemoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 3.调用逻辑处理函数
		l := logic.NewFirstdemoLogic(r.Context(), svcCtx)
		resp, err := l.Firstdemo(&req) // 4.调用返回的处理函数
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
