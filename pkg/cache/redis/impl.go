package redis

import (
	"context"
	"log"
	"time"

	"encoding/json"

	"github.com/Zhima-Mochi/linkZapURL/config"
	"github.com/Zhima-Mochi/linkZapURL/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type impl struct {
	client *redis.ClusterClient
	// client *redis.Client
	config *config.Redis
}

func NewRedis(config *config.Redis) (cache.Cache, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.Addrs,
		Password: config.Password,

		RouteByLatency: true,
		RouteRandomly:  true,
	})

	// client := redis.NewClient(&redis.Options{
	// 	Addr:     config.Addrs[0],
	// 	Password: config.Password,
	// })

	// Ping the primary
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis!")

	return &impl{
		client: client,
		config: config,
	}, nil
}

func (im *impl) Get(ctx context.Context, key string, result interface{}) error {
	val, err := im.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return cache.ErrNotFound
	} else if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), result)
	if err != nil {
		return err
	}

	return nil
}

func (im *impl) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return im.client.Set(ctx, key, jsonData, ttl).Err()
}

func (im *impl) Del(ctx context.Context, key string) error {
	return im.client.Del(ctx, key).Err()
}
