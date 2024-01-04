package message_queue

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

// RedisMessagePublish 发布消息到消息队列
func RedisMessagePublish(rdb *redis.Client, channel string, data interface{}) (int64, error) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	cnt, err := rdb.Publish(channel, marshal).Result()
	return cnt, err
}

// RedisMessagePSubscribeChannels 订阅消息队列
func RedisMessagePSubscribeChannels(rdb *redis.Client, channels ...string) <-chan *redis.Message {
	pubSub := rdb.PSubscribe(channels...)
	return pubSub.Channel()
}
