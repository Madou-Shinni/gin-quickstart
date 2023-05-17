package data

import (
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
)

type DemoRepo struct {
}

func (s *DemoRepo) Create(demo domain.Demo) error {
	return global.DB.Create(&demo).Error
}

func (s *DemoRepo) Delete(demo domain.Demo) error {
	return global.DB.Delete(&demo).Error
}

func (s *DemoRepo) DeleteByIds(ids request.Ids) error {
	return global.DB.Delete(&[]domain.Demo{}, ids.Ids).Error
}

func (s *DemoRepo) Update(demo map[string]interface{}) error {
	var columns []string
	for key := range demo {
		columns = append(columns, key)
	}
	if _, ok := demo["id"]; !ok {
		// 不存在id
		return errors.New(fmt.Sprintf("missing %s.id", "demo"))
	}
	model := domain.Demo{}
	model.ID = uint(demo["id"].(float64))
	return global.DB.Model(&model).Select(columns).Updates(&demo).Error
}

func (s *DemoRepo) Find(demo domain.Demo) (domain.Demo, error) {
	db := global.DB.Model(&domain.Demo{})
	// TODO：条件过滤

	res := db.First(&demo)

	return demo, res.Error
}

func (s *DemoRepo) List(page domain.PageDemoSearch) ([]domain.Demo, error) {
	var (
		demoList []domain.Demo
		err      error
	)
	// db
	db := global.DB.Model(&domain.Demo{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	// TODO：条件过滤

	err = db.Offset(offset).Limit(limit).Find(&demoList).Error

	return demoList, err
}

func (s *DemoRepo) Count() (int64, error) {
	var (
		count int64
		err   error
	)

	err = global.DB.Model(&domain.Demo{}).Count(&count).Error

	return count, err
}
