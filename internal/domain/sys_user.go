package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type SysUser struct {
	model.Model
	Account     string    `gorm:"size:255;unique;not null" json:"account"`   // 账号
	Password    string    `gorm:"size:255;not null" json:"password"`         // 密码
	NickName    string    `gorm:"size:255;not null" json:"nick_name"`        // 昵称
	DefaultRole uint      `gorm:"column:default_role" json:"default_role"`   // 当前角色
	Roles       []SysRole `gorm:"many2many:sys_user_sys_role;" json:"roles"` // 角色列表
}

type PageSysUserSearch struct {
	SysUser
	request.PageSearch
}

func (SysUser) TableName() string {
	return "sys_user"
}

type LoginReq struct {
	Account  string `gorm:"size:255;unique;not null" json:"account" binding:"required"` // 账号
	Password string `gorm:"size:255;not null" json:"password" binding:"required"`       // 密码
}
