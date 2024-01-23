//go:generate mockgen -destination=mock_cache.go -package=cache -source=cache.go
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
	Get(ctx context.Context, key string, result interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}
