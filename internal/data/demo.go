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

type DemoRepo struct {
}

func (s *DemoRepo) Create(ctx context.Context, demo domain.Demo) error {
	return global.DB.WithContext(ctx).Create(&demo).Error
}

func (s *DemoRepo) Delete(ctx context.Context, demo domain.Demo) error {
	return global.DB.WithContext(ctx).Delete(&demo).Error
}

func (s *DemoRepo) DeleteByIds(ctx context.Context, ids request.Ids) error {
	return global.DB.WithContext(ctx).Delete(&[]domain.Demo{}, ids.Ids).Error
}

func (s *DemoRepo) Update(ctx context.Context, demo map[string]interface{}) error {
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
	return global.DB.WithContext(ctx).Model(&model).Select(columns).Updates(&demo).Error
}

func (s *DemoRepo) Find(ctx context.Context, demo domain.Demo) (domain.Demo, error) {
	db := global.DB.Model(&domain.Demo{}).WithContext(ctx)
	// TODO：条件过滤

	res := db.First(&demo)

	return demo, res.Error
}

func (s *DemoRepo) List(ctx context.Context, page domain.PageDemoSearch) ([]domain.Demo, int64, error) {
	var (
		demoList []domain.Demo
		count    int64
		err      error
	)
	// db
	db := global.DB.Model(&domain.Demo{}).WithContext(ctx)
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	// TODO：条件过滤

	err = db.Count(&count).Offset(offset).Limit(limit).Find(&demoList).Error

	return demoList, count, err
}
