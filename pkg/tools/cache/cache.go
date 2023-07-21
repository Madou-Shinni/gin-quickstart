package cache

import (
	"time"
)

type Cache interface {
	Set(key string, value interface{}, expire time.Duration) error
	Get(key string) (interface{}, error)
	Del(key string) error
	DelByPrefix(key string) error
}
