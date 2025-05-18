package turtlesoup

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"turtle-soup/internal/logic/turtlesoup"
	"turtle-soup/internal/svc"
	types "turtle-soup/internal/types/turtlesoup"
)

func GetActivityToken(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetActivityTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := turtlesoup.NewGetActivityToken(r.Context(), svcCtx, r)
		resp, err := l.GetActivityToken(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
