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

func AddBookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookAdd
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := books.NewAddBookLogic(r.Context(), svcCtx)
		resp, err := l.AddBook(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
