package service

import (
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
)

func SmsSend(phoneNumber string, signName string, templateCode string, templateParams any) error {
	resp, err := global.SMS.Send(phoneNumber, signName, templateCode, templateParams)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return errors.New(fmt.Sprintf("SmsSend err text: %s", resp.Text))
	}

	return nil
}
