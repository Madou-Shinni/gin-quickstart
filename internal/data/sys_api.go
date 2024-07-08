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

type SysApiRepo struct {
}

func (s *SysApiRepo) Create(ctx context.Context, sysApi domain.SysApi) error {
	return global.DB.WithContext(ctx).Create(&sysApi).Error
}

func (s *SysApiRepo) Delete(ctx context.Context, sysApi domain.SysApi) error {
	return global.DB.WithContext(ctx).Delete(&sysApi).Error
}

func (s *SysApiRepo) DeleteByIds(ctx context.Context, ids request.Ids) error {
	return global.DB.WithContext(ctx).Delete(&[]domain.SysApi{}, ids.Ids).Error
}

func (s *SysApiRepo) Update(ctx context.Context, sysApi map[string]interface{}) error {
	var columns []string
	for key := range sysApi {
		columns = append(columns, key)
	}
	if _, ok := sysApi["id"]; !ok {
		// 不存在id
		return errors.New(fmt.Sprintf("missing %s.id", "sysApi"))
	}
	model := domain.SysApi{}
	model.ID = uint(sysApi["id"].(float64))
	return global.DB.WithContext(ctx).Model(&model).Select(columns).Updates(&sysApi).Error
}

func (s *SysApiRepo) Find(ctx context.Context, sysApi domain.SysApi) (domain.SysApi, error) {
	db := global.DB.WithContext(ctx).Model(&domain.SysApi{})
	// TODO：条件过滤

	res := db.First(&sysApi)

	return sysApi, res.Error
}

func (s *SysApiRepo) List(ctx context.Context, page domain.PageSysApiSearch) ([]domain.SysApi, int64, error) {
	var (
		sysApiList []domain.SysApi
		count      int64
		err        error
	)
	// db
	db := global.DB.WithContext(ctx).Model(&domain.SysApi{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	// TODO：条件过滤

	err = db.Count(&count).Offset(offset).Limit(limit).Find(&sysApiList).Error

	return sysApiList, count, err
}
