package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type RdbCache struct {
	rdb *redis.Client
}

func NewRdbCache(rdb *redis.Client) *RdbCache {
	return &RdbCache{rdb: rdb}
}

func (r *RdbCache) Set(key string, value interface{}, expire time.Duration) error {
	marshal, _ := json.Marshal(value)
	return r.rdb.SetNX(key, marshal, expire).Err()
}

func (r *RdbCache) Get(key string) (interface{}, error) {
	result, err := r.rdb.Get(key).Result()
	if err != nil {
		return nil, err
	}

	if result != "" {
		// 延长缓存的时间
		if r.rdb.TTL(key).Val() < 5 {
			r.rdb.Expire(key, time.Second*5)
		}
	}

	return result, nil
}

func (r *RdbCache) Del(key string) error {
	return r.rdb.Del(key).Err()
}

func (r *RdbCache) DelByPrefix(pre string) error {
	iter := r.rdb.Scan(0, pre+"*", 0).Iterator()
	for iter.Next() {
		if err := r.rdb.Del(iter.Val()).Err(); err != nil {
			return err
		}
	}
	return nil
}
