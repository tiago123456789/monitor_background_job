package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type CacheInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string)
}

type Cache struct {
	Client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{Client: client}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value string) error {
	err := c.Client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
