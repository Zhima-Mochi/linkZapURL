//go:generate mockgen -destination=mock_redirection.go -package=redirection -source=redirection.go
package redirection

import (
	"context"
	"errors"

	"github.com/Zhima-Mochi/linkZapURL/models"
)

var (
	ErrExpired = errors.New("redirection: expired")

	ErrNotFound = errors.New("redirection: not found")

	ErrInvalidCode = errors.New("redirection: invalid code")
)

type Redirection interface {
	Redirect(ctx context.Context, code string) (*models.URL, error)
}
