package sms

import (
	"encoding/json"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools"
	"time"
)

const defaultTimeout = time.Second

type ISms interface {
	Send(phoneNumber, signName, templateCode string, templateParams any) (sendResp, error)
}

type TemplateParams any

type sms struct {
	url            string
	timeOut        time.Duration
	headers        map[string]string
	phoneNumber    string
	signName       string
	templateCode   string
	templateParams TemplateParams
}

type sendResp struct {
	Code int         `json:"code"`
	Text string      `json:"text"`
	Data interface{} `json:"data"`
}

func NewSms(timeOut time.Duration, headers map[string]string, url string) *sms {
	s := &sms{url: url, timeOut: timeOut, headers: headers}
	if timeOut < time.Second {
		s.timeOut = defaultTimeout
	}
	return s
}

func (s *sms) Send(phoneNumber, signName, templateCode string, templateParams any) (sendResp, error) {
	s.signName = signName
	s.templateCode = templateCode
	s.phoneNumber = phoneNumber
	s.templateParams = templateParams
	return s.send()
}

func (s *sms) send() (sendResp, error) {
	var resp sendResp
	data := map[string]interface{}{
		"phone_number":   s.phoneNumber,
		"template_code":  s.templateCode,
		"sign_name":      s.signName,
		"template_param": s.templateParams,
	}

	body, err := tools.NewRequest(tools.POST, s.timeOut, s.url, data, s.headers)

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return resp, err
	}
	return resp, err
}
