package errors

import "github.com/pkg/errors"

var (
	ErrCacheNotFound        = errors.New("cache not found")
	ErrActiveTokenExpired   = errors.New("活动token已过期")
	ErrNoMoreAttempts       = errors.New("没有更多尝试次数了")
	ErrSessionTokenNotFound = errors.New("未找到会话token")
)
