package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type SysApi struct {
	model.Model
	Name   string `json:"name"`
	Method string `json:"method" gorm:"index:idx_method_path,unique"`
	Path   string `json:"path" gorm:"index:idx_method_path,unique"`
}

type PageSysApiSearch struct {
	SysApi
	request.PageSearch
}

func (SysApi) TableName() string {
	return "sys_api"
}
