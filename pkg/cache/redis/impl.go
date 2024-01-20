package redis

import (
	"context"
	"time"

	"encoding/json"

	"github.com/Zhima-Mochi/linkZapURL/config"
	"github.com/Zhima-Mochi/linkZapURL/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type impl struct {
	client *redis.ClusterClient
	config *config.Redis
}

func NewRedis(config *config.Redis) (cache.Cache, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: config.Addrs,
	})

	// Ping the primary
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &impl{
		client: client,
		config: config,
	}, nil
}

func (im *impl) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := im.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, cache.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (im *impl) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return im.client.Set(ctx, key, jsonData, ttl).Err()
}