package data

import (
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
)

type SysUserRepo struct {
}

func (s *SysUserRepo) Create(sysUser domain.SysUser) error {
	return global.DB.Create(&sysUser).Error
}

func (s *SysUserRepo) Delete(sysUser domain.SysUser) error {
	return global.DB.Delete(&sysUser).Error
}

func (s *SysUserRepo) DeleteByIds(ids request.Ids) error {
	return global.DB.Delete(&[]domain.SysUser{}, ids.Ids).Error
}

func (s *SysUserRepo) Update(sysUser map[string]interface{}) error {
	var columns []string
	for key := range sysUser {
		columns = append(columns, key)
	}
	if _, ok := sysUser["id"]; !ok {
		// 不存在id
		return errors.New(fmt.Sprintf("missing %s.id", "sysUser"))
	}
	model := domain.SysUser{}
	model.ID = uint(sysUser["id"].(float64))
	return global.DB.Model(&model).Select(columns).Updates(&sysUser).Error
}

func (s *SysUserRepo) Find(sysUser domain.SysUser) (domain.SysUser, error) {
	db := global.DB.Model(&domain.SysUser{})
	// TODO：条件过滤

	res := db.First(&sysUser)

	return sysUser, res.Error
}

func (s *SysUserRepo) List(page domain.PageSysUserSearch) ([]domain.SysUser, int64, error) {
	var (
		sysUserList []domain.SysUser
		count       int64
		err         error
	)
	// db
	db := global.DB.Model(&domain.SysUser{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	// TODO：条件过滤

	err = db.Count(&count).Offset(offset).Limit(limit).Find(&sysUserList).Error

	return sysUserList, count, err
}
