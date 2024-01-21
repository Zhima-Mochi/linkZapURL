package shortening

import (
	"context"
	"errors"

	"github.com/Zhima-Mochi/linkZapURL/models"
)

var (
	ErrEmptyURL = errors.New("shortening: empty url")

	ErrInvalidExpireAt = errors.New("shortening: invalid expire at")
)

type Shortening interface {
	Shorten(ctx context.Context, url string, expireAt int64) (*models.URL, error)
}
