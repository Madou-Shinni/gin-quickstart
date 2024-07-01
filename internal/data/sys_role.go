package data

import (
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
)

type SysRoleRepo struct {
}

func (s *SysRoleRepo) Create(sysRole domain.SysRole) error {
	return global.DB.Create(&sysRole).Error
}

func (s *SysRoleRepo) Delete(sysRole domain.SysRole) error {
	return global.DB.Delete(&sysRole).Error
}

func (s *SysRoleRepo) DeleteByIds(ids request.Ids) error {
	return global.DB.Delete(&[]domain.SysRole{}, ids.Ids).Error
}

func (s *SysRoleRepo) Update(sysRole map[string]interface{}) error {
	var columns []string
	for key := range sysRole {
		columns = append(columns, key)
	}
	if _, ok := sysRole["id"]; !ok {
		// 不存在id
		return errors.New(fmt.Sprintf("missing %s.id", "sysRole"))
	}
	model := domain.SysRole{}
	model.ID = uint(sysRole["id"].(float64))
	return global.DB.Model(&model).Select(columns).Updates(&sysRole).Error
}

func (s *SysRoleRepo) Find(sysRole domain.SysRole) (domain.SysRole, error) {
	db := global.DB.Model(&domain.SysRole{})
	// TODO：条件过滤

	res := db.First(&sysRole)

	return sysRole, res.Error
}

func (s *SysRoleRepo) List(page domain.PageSysRoleSearch) ([]domain.SysRole, int64, error) {
	var (
		sysRoleList []domain.SysRole
		count       int64
		err         error
	)
	// db
	db := global.DB.Model(&domain.SysRole{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	// TODO：条件过滤

	err = db.Count(&count).Offset(offset).Limit(limit).Find(&sysRoleList).Error

	return sysRoleList, count, err
}
