package task

import (
	"encoding/json"
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/go-logger"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"time"
)

const defaultMaxRetry = 3

type taskV2Opt func(client *TaskV2)

func WithMaxRetry(maxRetry int) taskV2Opt {
	return func(client *TaskV2) {
		client.maxRetry = maxRetry
	}
}

type TaskV2 struct {
	maxRetry int
	errCh    chan error
	*asynq.Scheduler
}

func NewTaskV2(s *asynq.Scheduler, opts ...taskV2Opt) *TaskV2 {
	errCh := make(chan error, 1)
	go func() {
		for {
			errHandle(<-errCh)
		}
	}()

	taskV2 := &TaskV2{Scheduler: s, errCh: errCh}

	for _, opt := range opts {
		opt(taskV2)
	}

	return taskV2
}

func (t *TaskV2) NewTask(queue string, payload any, opts ...asynq.Option) *asynq.Task {
	marshal, err := json.Marshal(payload)
	if err != nil {
		t.errCh <- err
	}

	opts = append([]asynq.Option{asynq.MaxRetry(t.maxRetry)}, opts...)

	return asynq.NewTask(queue, marshal, opts...)
}

func (t *TaskV2) Register(spec string, task *asynq.Task) string {
	entryID, err := t.Scheduler.Register(spec, task)
	if err != nil {
		t.errCh <- err
	}
	return entryID
}

// V2Init Scheduler
/*
cron 表达式 https://tooltt.com/crontab/c/35.html
*/
func V2Init() {
	config := conf.Conf.AsynqConfig
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// 周期性任务
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		}, &asynq.SchedulerOpts{
			Location:        loc,
			PostEnqueueFunc: handleEnqueueError,
		})

	v2 := NewTaskV2(scheduler, WithMaxRetry(defaultMaxRetry))

	v2.Register("@every 30m", v2.NewTask(constants.TaskTest, nil)) // 每隔30分钟同步一次

	if err := scheduler.Run(); err != nil {
		logger.Error("could not run scheduler", zap.Error(err))
	}
}

func handleEnqueueError(task *asynq.TaskInfo, err error) {
	// 你的错误处理逻辑
	if err != nil {
		logger.Error("定时任务入队失败", zap.Error(err), zap.Any("taskInfo", task))
	}
}
