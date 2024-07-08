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

type FileRepo struct {
}

func (s *FileRepo) Create(ctx context.Context, file domain.File) error {
	return global.DB.WithContext(ctx).Create(&file).Error
}

func (s *FileRepo) Delete(ctx context.Context, file domain.File) error {
	return global.DB.WithContext(ctx).Delete(&file).Error
}

func (s *FileRepo) DeleteByIds(ctx context.Context, ids request.Ids) error {
	return global.DB.WithContext(ctx).Delete(&[]domain.File{}, ids.Ids).Error
}

func (s *FileRepo) Update(ctx context.Context, file map[string]interface{}) error {
	var columns []string
	for key := range file {
		columns = append(columns, key)
	}
	if _, ok := file["id"]; !ok {
		// 不存在id
		return errors.New(fmt.Sprintf("missing %s.id", "file"))
	}
	model := domain.File{}
	model.ID = int64(file["id"].(float64))
	return global.DB.WithContext(ctx).Model(&model).Select(columns).Updates(&file).Error
}

func (s *FileRepo) Find(ctx context.Context, file domain.File) (domain.File, error) {
	db := global.DB.WithContext(ctx).Model(&domain.File{})
	if file.FileMd5 != "" {
		db = db.Where("file_md5 = ?", file.FileMd5)
	}

	res := db.First(&file)

	return file, res.Error
}

func (s *FileRepo) List(ctx context.Context, page domain.PageFileSearch) ([]domain.File, int64, error) {
	var (
		fileList []domain.File
		count    int64
		err      error
	)
	// db
	db := global.DB.WithContext(ctx).Model(&domain.File{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	err = db.Count(&count).Offset(offset).Limit(limit).Find(&fileList).Error

	return fileList, count, err
}
