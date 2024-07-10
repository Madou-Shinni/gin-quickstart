package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type SysRole struct {
	model.Model
	ParentID uint      `gorm:"column:parent_id" json:"parent_id"`
	RoleName string    `gorm:"column:role_name" json:"role_name"`
	Menus    []SysMenu `gorm:"many2many:sys_role_sys_menu;" json:"menus"` // 菜单列表
	Children []SysRole `gorm:"-" json:"children"`                         // 角色列表
}

type PageSysRoleSearch struct {
	SysRole
	request.PageSearch
}

func (SysRole) TableName() string {
	return "sys_role"
}
