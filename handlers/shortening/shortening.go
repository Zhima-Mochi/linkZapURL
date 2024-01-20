package shortening

import (
	"context"
	"errors"
)

var (
	ErrEmptyURL = errors.New("shortening: empty url")
)

type Shortening interface {
	Shorten(ctx context.Context, url string) (string, error)
}
