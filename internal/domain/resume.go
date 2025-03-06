package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type Resume struct {
	model.Model
}

type PageResumeSearch struct {
	Resume
	request.PageSearch
}

func (Resume) TableName() string {
	return "resume"
}
