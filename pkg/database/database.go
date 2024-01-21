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
	Get(ctx context.Context, table string, key int64) (interface{}, error)
	Set(ctx context.Context, table string, key int64, value interface{}) error
}
