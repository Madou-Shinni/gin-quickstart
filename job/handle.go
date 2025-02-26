package job

import (
	"context"
	"encoding/json"
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/go-logger"
	"github.com/hibiken/asynq"
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
	var payload domain.DataImport
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	switch payload.Category {
	case constants.DataImportCategoryDemo:
		err = importDemo(ctx, payload)
	}

	return err
}
