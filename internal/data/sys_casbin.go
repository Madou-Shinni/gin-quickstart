package data

import (
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
)

type SysCasbinRepo struct {
}

func (s *SysCasbinRepo) Create(sysCasbin domain.SysCasbin) error {
	return global.DB.Create(&sysCasbin).Error
}

func (s *SysCasbinRepo) Delete(sysCasbin domain.SysCasbin) error {
	return global.DB.Delete(&sysCasbin).Error
}

func (s *SysCasbinRepo) DeleteByIds(ids request.Ids) error {
	return global.DB.Delete(&[]domain.SysCasbin{}, ids.Ids).Error
}

func (s *SysCasbinRepo) Update(sysCasbin map[string]interface{}) error {
	var columns []string
	for key := range sysCasbin {
		columns = append(columns, key)
	}
	if _, ok := sysCasbin["id"]; !ok {
		// 不存在id
		return errors.New(fmt.Sprintf("missing %s.id", "sysCasbin"))
	}
	model := domain.SysCasbin{}
	return global.DB.Model(&model).Select(columns).Updates(&sysCasbin).Error
}

func (s *SysCasbinRepo) Find(sysCasbin domain.SysCasbin) (domain.SysCasbin, error) {
	db := global.DB.Model(&domain.SysCasbin{})
	// TODO：条件过滤

	res := db.First(&sysCasbin)

	return sysCasbin, res.Error
}

func (s *SysCasbinRepo) List(page domain.PageSysCasbinSearch) ([]domain.SysCasbin, int64, error) {
	var (
		sysCasbinList []domain.SysCasbin
		count         int64
		err           error
	)
	// db
	db := global.DB.Model(&domain.SysCasbin{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	// TODO：条件过滤

	err = db.Count(&count).Offset(offset).Limit(limit).Find(&sysCasbinList).Error

	return sysCasbinList, count, err
}
