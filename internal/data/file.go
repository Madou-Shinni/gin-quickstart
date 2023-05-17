package data

import (
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
)

type FileRepo struct {
}

func (s *FileRepo) Create(file domain.File) error {
	return global.DB.Create(&file).Error
}

func (s *FileRepo) Delete(file domain.File) error {
	return global.DB.Delete(&file).Error
}

func (s *FileRepo) DeleteByIds(ids request.Ids) error {
	return global.DB.Delete(&[]domain.File{}, ids.Ids).Error
}

func (s *FileRepo) Update(file map[string]interface{}) error {
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
	return global.DB.Model(&model).Select(columns).Updates(&file).Error
}

func (s *FileRepo) Find(file domain.File) (domain.File, error) {
	db := global.DB.Model(&domain.File{})
	if file.FileMd5 != "" {
		db = db.Where("file_md5 = ?", file.FileMd5)
	}

	res := db.First(&file)

	return file, res.Error
}

func (s *FileRepo) List(page domain.PageFileSearch) ([]domain.File, error) {
	var (
		fileList []domain.File
		err      error
	)
	// db
	db := global.DB.Model(&domain.File{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	err = db.Offset(offset).Limit(limit).Find(&fileList).Error

	return fileList, err
}

func (s *FileRepo) Count() (int64, error) {
	var (
		count int64
		err   error
	)

	err = global.DB.Model(&domain.File{}).Count(&count).Error

	return count, err
}
