package job

import (
	"github.com/Madou-Shinni/gin-quickstart/constants"
	_ "github.com/Madou-Shinni/gin-quickstart/initialize"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandle(t *testing.T) {
	var payload domain.Sms
	payload.PhoneNumber = "1888888888"
	err := global.Producer.NewTask(constants.QueueSms, payload)
	assert.Error(t, err)
}
