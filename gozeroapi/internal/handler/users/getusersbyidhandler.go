// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozeroapi/internal/logic/users"
	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"
)

func GetUsersByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserQuestById
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := users.NewGetUsersByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetUsersById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
