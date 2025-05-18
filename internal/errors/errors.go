package errors

import "github.com/pkg/errors"

var (
	ErrCacheNotFound = errors.New("cache not found")
)
