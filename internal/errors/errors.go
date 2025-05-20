package errors

import "github.com/pkg/errors"

var (
	ErrCacheNotFound      = errors.New("cache not found")
	ErrActiveTokenExpired = errors.New("active token expired")
	ErrNoMoreAttempts     = errors.New("没有更多尝试次数了")
)
