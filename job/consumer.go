package job

import (
	"context"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/go-logger"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"log"
)

func RunConsumer() {
	asynqConfig := conf.Conf.AsynqConfig
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     asynqConfig.Addr,
			Password: asynqConfig.Password,
			DB:       asynqConfig.DB,
		},
		asynq.Config{
			// 每个进程并发执行的worker数量
			Concurrency:  20,
			ErrorHandler: asynq.ErrorHandlerFunc(errHandlerFunc),
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(constants.QueueSms, handleSmsSend)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func errHandlerFunc(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		id, ok := asynq.GetTaskID(ctx)
		if !ok {
			id = "unknown"
		}
		err = fmt.Errorf("retry exhausted for job %s, task id [%s] err: %w", task.Type(), id, err)
	}
	//errorReportingService. Notify(err)
	logger.Error("消费异常", zap.Error(err))
}
