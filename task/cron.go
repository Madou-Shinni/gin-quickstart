package task

import (
	"github.com/Madou-Shinni/go-logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var c *cron.Cron

func V1Init() {
	// Seconds field, optional
	c = cron.New(cron.WithSeconds())

	c.Start()
}

// 只能监测启动时的错误
func errHandle(err error) {
	logger.Error("添加定时任务发生错误", zap.Error(err))
}
