package turtlesoup

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"turtle-soup/internal/logic/turtlesoup"
	"turtle-soup/internal/svc"
	types "turtle-soup/internal/types/turtlesoup"
)

func GetSessionToken(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSessionTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := turtlesoup.NewGetSessionToken(r.Context(), svcCtx, r)
		resp, err := l.GetSessionToken(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
