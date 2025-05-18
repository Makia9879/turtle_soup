package turtlesoup

import (
	"net/http"

	"context"

	"turtle-soup/internal/svc"
	types "turtle-soup/internal/types/turtlesoup"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionToken struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewGetSessionToken(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *GetSessionToken {
	return &GetSessionToken{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *GetSessionToken) GetSessionToken(req *types.GetSessionTokenRequest) (resp *types.GetSessionTokenResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
