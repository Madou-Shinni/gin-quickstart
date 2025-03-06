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

type ResumeRepo struct {
}

func (s *ResumeRepo) Create(ctx context.Context, resume *domain.Resume) error {
	return global.DB.WithContext(ctx).Create(&resume).Error
}

func (s *ResumeRepo) Delete(ctx context.Context, resume domain.Resume) error {
	return global.DB.WithContext(ctx).Delete(&resume).Error
}

func (s *ResumeRepo) DeleteByIds(ctx context.Context, ids request.Ids) error {
	return global.DB.WithContext(ctx).Delete(&[]domain.Resume{}, ids.Ids).Error
}

func (s *ResumeRepo) Update(ctx context.Context, resume domain.Resume) error {
	if resume.ID == 0 {
		return errors.New(fmt.Sprintf("missing %s.id", "resume"))
	}
	return global.DB.WithContext(ctx).Model(&domain.Resume{Model: model.Model{ID: resume.ID}}).Updates(&resume).Error
}

func (s *ResumeRepo) Find(ctx context.Context, resume domain.Resume) (domain.Resume, error) {
	db := global.DB.WithContext(ctx).Model(&domain.Resume{})
	// TODO：条件过滤

	res := db.First(&resume)

	return resume, res.Error
}

func (s *ResumeRepo) List(ctx context.Context, page domain.PageResumeSearch) ([]domain.Resume, int64, error) {
	var (
		resumeList []domain.Resume
		count      int64
		err        error
	)
	// db
	db := global.DB.WithContext(ctx).Model(&domain.Resume{})

	// TODO：条件过滤

	err = db.Count(&count).Scopes(scopes.Paginate(page.PageSearch), scopes.OrderBy(page.OrderBy)).Find(&resumeList).Error

	return resumeList, count, err
}
