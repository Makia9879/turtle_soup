package version

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"turtle-soup/internal/logic/version"
	"turtle-soup/internal/svc"
	types "turtle-soup/internal/types/version"
)

func Get(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := version.NewGet(r.Context(), svcCtx, r)
		resp, err := l.Get(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
