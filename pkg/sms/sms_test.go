package sms

import (
	"fmt"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	params := struct {
		Operation string `json:"operation"`
		Msg       string `json:"msg"`
	}{
		Operation: "入住",
		Msg:       "这是一条测试发送",
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Token %s", "71f98c5037ad49adbd16cf93610d8cdd"),
	}
	url := "https://sms.iqusong.com" + "/v1/sms"
	smsSender := NewSms(1*time.Second, headers, url)
	resp, err := smsSender.Send("17867897925", "才得利", "SMS_11111283", params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}
