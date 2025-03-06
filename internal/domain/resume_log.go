package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type ResumeLog struct {
	model.Model
	RID uint `json:"rid"` // 简历
	IID uint `json:"iid"` // 面试点
}

type PageResumeLogSearch struct {
	ResumeLog
	request.PageSearch
}

func (ResumeLog) TableName() string {
	return "resume_log"
}
