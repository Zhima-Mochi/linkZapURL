package database

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("database: key not found")

	ErrCollision = errors.New("database: collision")
)

type Database interface {
	Get(ctx context.Context, table, key string) (interface{}, error)
	Set(ctx context.Context, table, key string, value interface{}) error
}
