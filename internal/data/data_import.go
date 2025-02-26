package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/scopes"
)

type DataImportRepo struct {
}

func (s *DataImportRepo) Create(ctx context.Context, dataImport *domain.DataImport) error {
	return global.DB.WithContext(ctx).Create(&dataImport).Error
}

func (s *DataImportRepo) Delete(ctx context.Context, dataImport domain.DataImport) error {
	return global.DB.WithContext(ctx).Delete(&dataImport).Error
}

func (s *DataImportRepo) DeleteByIds(ctx context.Context, ids request.Ids) error {
	return global.DB.WithContext(ctx).Delete(&[]domain.DataImport{}, ids.Ids).Error
}

func (s *DataImportRepo) Update(ctx context.Context, dataImport domain.DataImport) error {
	if dataImport.ID == 0 {
		return errors.New(fmt.Sprintf("missing %s.id", "dataImport"))
	}
	return global.DB.WithContext(ctx).Model(&domain.DataImport{Model: model.Model{ID: dataImport.ID}}).Updates(&dataImport).Error
}

func (s *DataImportRepo) Find(ctx context.Context, dataImport domain.DataImport) (domain.DataImport, error) {
	db := global.DB.WithContext(ctx).Model(&domain.DataImport{})
	// TODO：条件过滤

	res := db.First(&dataImport)

	return dataImport, res.Error
}

func (s *DataImportRepo) List(ctx context.Context, page domain.PageDataImportSearch) ([]domain.DataImport, int64, error) {
	var (
		dataImportList []domain.DataImport
		count          int64
		err            error
	)
	// db
	db := global.DB.WithContext(ctx).Model(&domain.DataImport{})

	// TODO：条件过滤

	err = db.Count(&count).Scopes(scopes.Paginate(page.PageSearch)).Find(&dataImportList).Error

	return dataImportList, count, err
}
