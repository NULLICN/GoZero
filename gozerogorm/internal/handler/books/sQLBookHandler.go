// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package books

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerogorm/internal/logic/books"
	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"
)

func SQLBookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookSQL
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := books.NewSQLBookLogic(r.Context(), svcCtx)
		resp, err := l.SQLBook(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
