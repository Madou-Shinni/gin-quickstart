package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type SysCasbin struct {
}

type PageSysCasbinSearch struct {
	SysCasbin
	request.PageSearch
}

func (SysCasbin) TableName() string {
	return "sys_casbin"
}

type UserRolesReq struct {
	UserID uint   `json:"user_id"`
	Roles  []uint `json:"roles"`
}

type RoleRolesReq struct {
	Role  uint   `json:"role"`
	Roles []uint `json:"roles"`
}

type RolePermissionsReq struct {
	Role        uint              `json:"role"`
	Permissions map[string]string `json:"permissions"`
}
