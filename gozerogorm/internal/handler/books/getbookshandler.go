// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package books

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerogorm/internal/logic/books"
	"gozerogorm/internal/svc"
)

func GetbooksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := books.NewGetbooksLogic(r.Context(), svcCtx)
		resp, err := l.Getbooks()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
