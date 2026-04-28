// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lesson

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerogorm/internal/logic/lesson"
	"gozerogorm/internal/svc"
)

func GetLessonListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := lesson.NewGetLessonListLogic(r.Context(), svcCtx)
		resp, err := l.GetLessonList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
