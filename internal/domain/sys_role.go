package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type SysRole struct {
	model.Model
	ParentID uint   `gorm:"column:parent_id"`
	RoleName string `gorm:"column:role_name"`
}

type PageSysRoleSearch struct {
	SysRole
	request.PageSearch
}

func (SysRole) TableName() string {
	return "sys_role"
}
