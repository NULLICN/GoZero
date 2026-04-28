// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package student

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerogorm/internal/logic/student"
	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"
)

func GetOneStudentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StudentIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := student.NewGetOneStudentLogic(r.Context(), svcCtx)
		resp, err := l.GetOneStudent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
