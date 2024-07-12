package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type SysApi struct {
	ID        uint             `gorm:"primarykey" json:"id" form:"id" uri:"id"`
	CreatedAt *model.LocalTime `json:"createdAt" form:"createdAt" swaggerignore:"true"`
	UpdatedAt *model.LocalTime `json:"updatedAt" form:"updatedAt" swaggerignore:"true"`
	Name      string           `json:"name" gorm:"name"`
	Method    string           `json:"method" gorm:"index:idx_method_path,unique"`
	Path      string           `json:"path" gorm:"index:idx_method_path,unique"`
}

type PageSysApiSearch struct {
	SysApi
	request.PageSearch
}

func (SysApi) TableName() string {
	return "sys_api"
}
