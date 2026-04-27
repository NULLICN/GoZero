// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package bookmetadata

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerogorm/internal/logic/bookmetadata"
	"gozerogorm/internal/svc"
)

func GetBookAndMetadataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := bookmetadata.NewGetBookAndMetadataLogic(r.Context(), svcCtx)
		resp, err := l.GetBookAndMetadata()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
