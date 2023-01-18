package data

import (
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
)

type DemoRepo struct {
}

func NewDemoRepo() service.DemoRepo {
	return &DemoRepo{}
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

func (s *DemoRepo) Update(demo domain.Demo) error {
	return global.DB.Updates(&demo).Error
}

func (s *DemoRepo) Find(demo domain.Demo) (domain.Demo, error) {
	db := global.DB.Model(&domain.Demo{})
	if demo.Did != 0 {
		db = db.Where("cid = ?", demo.Did)
	}

	if demo.Uid != 0 {
		db = db.Where("uid = ?", demo.Uid)
	}

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

	if page.Did != 0 {
		db = db.Where("cid = ?", page.Did)
	}

	if page.Uid != 0 {
		db = db.Where("uid = ?", page.Uid)
	}

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
