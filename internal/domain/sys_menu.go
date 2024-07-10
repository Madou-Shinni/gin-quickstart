package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type SysMenu struct {
	model.Model
	Name        string    `gorm:"size:255;unique;not null" json:"name"` // 菜单名称
	Icon        string    `gorm:"size:255;not null" json:"icon"`        // 图标
	ParentID    uint      `gorm:"default:0" json:"parent_id"`           // 上层菜单
	Description string    `gorm:"size:255;" json:"description"`         // 描述
	Children    []SysMenu `gorm:"-" json:"children"`
}

type PageSysMenuSearch struct {
	SysMenu
	request.PageSearch
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
