package database

import "context"

type Database interface {
	Get(ctx context.Context, table, key string) (interface{}, error)
	Set(ctx context.Context, table, key string, value interface{}) error
}
