package turtlesoup

import (
	"net/http"

	"context"

	"turtle-soup/internal/svc"
	types "turtle-soup/internal/types/turtlesoup"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitAnswer struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewSubmitAnswer(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *SubmitAnswer {
	return &SubmitAnswer{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *SubmitAnswer) SubmitAnswer(req *types.SubmitAnswerRequest) (resp *types.SubmitAnswerResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
