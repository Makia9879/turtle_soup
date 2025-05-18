package turtlesoup

import (
	"github.com/zeromicro/go-zero/core/contextx"
	"net/http"
	"time"

	"context"

	"turtle-soup/internal/svc"
	types "turtle-soup/internal/types/turtlesoup"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetActivityToken struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewGetActivityToken(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *GetActivityToken {
	return &GetActivityToken{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *GetActivityToken) GetActivityToken(req *types.GetActivityTokenRequest) (resp *types.GetActivityTokenResponse, err error) {
	// 1. 验证请求参数
	//if req.ActivityID == "" {
	//	return nil, types.ErrInvalidActivityID
	//}
	valCtx := contextx.ValueOnlyFrom(l.ctx)

	// 2. 生成活动令牌
	token := uuid.New().String()
	e := time.Duration(l.svcCtx.MustGetConfig().ActiveTokenExpire) * time.Second
	expireTime := time.Now().Add(e)

	// 3. 存储令牌信息 (这里需要根据实际存储方式实现)
	// 例如使用Redis: err = l.svcCtx.Redis.SetexCtx(l.ctx, "activity_token:"+token, req.ActivityID, 86400)
	// 这里只是示例，实际实现需要根据您的存储方案调整
	err = l.svcCtx.
		RedisConn.
		SetexCtx(valCtx, "activity_token:"+token, time.Now().String(), l.svcCtx.MustGetConfig().ActiveTokenExpire)
	if err != nil {
		l.Logger.Errorf("[GetActivityToken] error: %v", err)
		return nil, err
	}

	// 4. 返回响应
	return &types.GetActivityTokenResponse{
		Token:      token,
		ExpireTime: expireTime.Unix(),
	}, nil
}
