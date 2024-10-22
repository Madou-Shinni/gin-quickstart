package message_queue

import (
	"encoding/json"
	"github.com/Madou-Shinni/go-logger"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

const defaultMaxRetry = 0

type asynqOpt func(client *AsynqClient)

type AsynqClient struct {
	*asynq.Client
	maxRetry int
}

func WithMaxRetry(maxRetry int) asynqOpt {
	return func(client *AsynqClient) {
		client.maxRetry = maxRetry
	}
}

func NewAsynqClient(addr string, pwd string, db int, opts ...asynqOpt) *AsynqClient {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("new asynq client panic", zap.Error(err.(error)))
		}
	}()
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	asynqClient := &AsynqClient{Client: client, maxRetry: defaultMaxRetry}

	for _, f := range opts {
		f(asynqClient)
	}

	return asynqClient
}

func (c *AsynqClient) NewTask(typename string, payload any, opts ...asynq.Option) error {
	marshal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(typename, marshal)

	opts = append([]asynq.Option{asynq.MaxRetry(c.maxRetry)}, opts...)

	_, err = c.Enqueue(task, opts...)
	if err != nil {
		return err
	}
	return nil
}
