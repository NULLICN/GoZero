// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package student

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerogorm/internal/logic/student"
	"gozerogorm/internal/svc"
)

func GetStudentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := student.NewGetStudentListLogic(r.Context(), svcCtx)
		resp, err := l.GetStudentList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
