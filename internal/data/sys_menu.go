package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
)

type SysMenuRepo struct {
}

func (s *SysMenuRepo) Create(ctx context.Context, sysMenu domain.SysMenu) error {
	return global.DB.WithContext(ctx).Create(&sysMenu).Error
}

func (s *SysMenuRepo) Delete(ctx context.Context, sysMenu domain.SysMenu) error {
	return global.DB.WithContext(ctx).Delete(&sysMenu).Error
}

func (s *SysMenuRepo) DeleteByIds(ctx context.Context, ids request.Ids) error {
	return global.DB.WithContext(ctx).Delete(&[]domain.SysMenu{}, ids.Ids).Error
}

func (s *SysMenuRepo) Update(ctx context.Context, sysMenu map[string]interface{}) error {
	var columns []string
	for key := range sysMenu {
		columns = append(columns, key)
	}
	if _, ok := sysMenu["id"]; !ok {
		// 不存在id
		return errors.New(fmt.Sprintf("missing %s.id", "sysMenu"))
	}
	model := domain.SysMenu{}
	model.ID = uint(sysMenu["id"].(float64))
	return global.DB.WithContext(ctx).Model(&model).Select(columns).Updates(&sysMenu).Error
}

func (s *SysMenuRepo) Find(ctx context.Context, sysMenu domain.SysMenu) (domain.SysMenu, error) {
	db := global.DB.WithContext(ctx).Model(&domain.SysMenu{})
	// TODO：条件过滤

	res := db.First(&sysMenu)

	return sysMenu, res.Error
}

func (s *SysMenuRepo) List(ctx context.Context, page domain.PageSysMenuSearch) ([]domain.SysMenu, int64, error) {
	var (
		sysMenuList []domain.SysMenu
		count       int64
		err         error
	)
	// db
	db := global.DB.WithContext(ctx).Model(&domain.SysMenu{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	// TODO：条件过滤

	err = db.Where("parent_id = ?", 0).Count(&count).Offset(offset).Limit(limit).Find(&sysMenuList).Error

	return sysMenuList, count, err
}
