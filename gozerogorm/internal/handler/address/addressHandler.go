// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerogorm/internal/logic/address"
	"gozerogorm/internal/svc"
)

func AddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewAddressLogic(r.Context(), svcCtx)
		resp, err := l.Address()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
