package message_queue

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

type Message struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func TestRedisMessagePublish(t *testing.T) {
	// redis连接初始化
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:16379",
		Password: "",
		DB:       0,
		PoolSize: 8,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		t.Fatalf("redis ping err: %v", err)
		return
	}

	defer rdb.Close()

	for i := 0; i < 10; i++ {
		cnt, err := RedisMessagePublish[Message](rdb, "test", Message{
			ID:   uint64(i),
			Name: fmt.Sprintf("name_%v", i),
		})
		if err != nil {
			t.Fatalf("redis publish err: %v", err)
			return
		}
		t.Logf("publish success count:%d : %v", cnt, fmt.Sprintf("name_%v", i))
	}
}

func TestRedisMessagePSubscribe(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:16379",
		Password: "",
		DB:       0,
		PoolSize: 8,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		t.Fatalf("redis ping err: %v", err)
		return
	}

	defer rdb.Close()

	subscribe := RedisMessagePSubscribeChannels(rdb, "test", "pubsub")

	for {
		message, ok := <-subscribe
		if !ok {
			t.Fatalf("redis subscribe false")
			return
		}
		switch message.Channel {
		case "test":
			t.Log("test channel message: ", message.Payload)
		case "pubsub":
			t.Log("pubsub channel message: ", message.Payload)
		}
	}
}
