// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozeroapi/internal/logic/users"
	"gozeroapi/internal/svc"
)

func GetUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := users.NewGetUsersLogic(r.Context(), svcCtx)
		resp, err := l.GetUsers()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
