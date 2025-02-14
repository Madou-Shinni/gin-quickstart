package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	GetAndRenewal(ctx context.Context, key string) (interface{}, error)
	Del(ctx context.Context, key string) error
	DelByPrefix(ctx context.Context, key string) error
}
