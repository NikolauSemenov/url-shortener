package cache

import (
	"context"
	"log"
	"time"
	"url-shortener/internal/config"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
	Invalidate(key string) error
}

type RedisCache struct {
	redis *redis.Client
	ctx   context.Context
}

func NewCache(cfg *config.CacheConfig) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.CacheDsn,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis ping error: %v", err)
		return nil, err
	}
	return &RedisCache{redis: rdb, ctx: context.Background()}, nil
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.redis.Get(r.ctx, key).Result()
}

func (r *RedisCache) Set(key string, value string, ttl time.Duration) error {
	return r.redis.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisCache) Invalidate(key string) error {
	return r.redis.Del(r.ctx, key).Err()
}
