package initialize

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/sms"
	"time"
)

func init() {
	smsConfig := conf.Conf.SMSConfig
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Token %s", smsConfig.SmsToken),
	}
	url := fmt.Sprint(smsConfig.SmsServer, smsConfig.SmsSendPath)
	global.SMS = sms.NewSms(time.Second*10, headers, url)
}
