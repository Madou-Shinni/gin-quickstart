package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type RdbCache struct {
	rdb *redis.Client
}

func NewRdbCache(rdb *redis.Client) *RdbCache {
	return &RdbCache{rdb: rdb}
}

func (r *RdbCache) Set(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	marshal, _ := json.Marshal(value)
	return r.rdb.SetNX(ctx, key, marshal, expire).Err()
}

func (r *RdbCache) Get(ctx context.Context, key string) (interface{}, error) {
	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAndRenewal 获取并且续期
func (r *RdbCache) GetAndRenewal(ctx context.Context, key string) (interface{}, error) {
	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if result != "" {
		// 延长缓存的时间
		val := r.rdb.TTL(ctx, key).Val()
		if val < 5*time.Second {
			r.rdb.Expire(ctx, key, time.Second*20)
		}
	}

	return result, nil
}

func (r *RdbCache) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *RdbCache) DelByPrefix(ctx context.Context, pre string) error {
	iter := r.rdb.Scan(ctx, 0, pre+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := r.rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return nil
}
