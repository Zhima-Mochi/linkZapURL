package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("cache: key not found")
)

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}
