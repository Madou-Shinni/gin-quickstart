package job

import (
	"context"
	"encoding/json"
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/go-logger"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func handleSmsSend(ctx context.Context, task *asynq.Task) error {
	var payload domain.Sms
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	return service.SmsSend(payload.PhoneNumber, payload.SignName, payload.TemplateCode, payload.TemplateParams)
}

func handleTaskTest(ctx context.Context, task *asynq.Task) error {
	// 这里可以添加任务测试的逻辑
	logger.Info("handleTaskTest")
	return nil
}

func handleImportData(ctx context.Context, task *asynq.Task) error {
	// 捕获异常 导入数据错误不进入asynq存档,避免存档负载过大
	defer func() {
		if err := recover(); err != nil {
			logger.Error("导入数据异常", zap.Any("err", err))
		}
	}()

	var payload domain.DataImport
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	switch payload.Category {
	case constants.DataImportCategoryDemo:
		err = importDemo(ctx, payload)
	}

	if err != nil {
		panic(err)
	}

	return err
}
