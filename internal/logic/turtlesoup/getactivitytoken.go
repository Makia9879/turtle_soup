package turtlesoup

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/contextx"
	"net/http"
	"strings"
	"time"
	"turtle-soup/internal/constant"

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
	valCtx := contextx.ValueOnlyFrom(l.ctx)
	// 1. 检查是否还有生效的 activeToken
	scanRes, _, err := l.svcCtx.
		RedisConn.
		ScanCtx(valCtx, 0, fmt.Sprintf("%s*", constant.CacheKeyActivityToken), 1)
	if err != nil {
		l.Logger.Errorf("[GetActivityToken] ScanCtx() error: %v", err)
		return nil, err
	}
	if len(scanRes) > 0 {
		scanResKey := scanRes[0]
		ttlRes, err := l.svcCtx.RedisConn.TtlCtx(valCtx, scanResKey)
		if err != nil {
			l.Logger.Errorf("[GetActivityToken] TtlCtx() error: %v", err)
			return nil, err
		}
		return &types.GetActivityTokenResponse{
			Token:      strings.ReplaceAll(scanResKey, constant.CacheKeyActivityToken, ""),
			ExpireTime: time.Now().Add(time.Duration(ttlRes) * time.Second).Unix(),
		}, nil
	}

	// 2. 生成活动令牌
	token := uuid.New().String()
	e := time.Duration(l.svcCtx.MustGetConfig().ActiveTokenExpire) * time.Second
	expireTime := time.Now().Add(e)

	// 3. 存储令牌信息 (这里需要根据实际存储方式实现)
	// 例如使用Redis: err = l.svcCtx.Redis.SetexCtx(l.ctx, "activity_token:"+token, req.ActivityID, 86400)
	// 这里只是示例，实际实现需要根据您的存储方案调整
	err = l.svcCtx.
		RedisConn.
		SetexCtx(valCtx, constant.GetActivityToken(token), time.Now().String(), l.svcCtx.MustGetConfig().ActiveTokenExpire)
	if err != nil {
		l.Logger.Errorf("[GetActivityToken] SetexCtx() error: %v", err)
		return nil, err
	}

	// 4. 返回响应
	return &types.GetActivityTokenResponse{
		Token:      token,
		ExpireTime: expireTime.Unix(),
	}, nil
}
