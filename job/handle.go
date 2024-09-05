package job

import (
	"context"
	"encoding/json"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
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
