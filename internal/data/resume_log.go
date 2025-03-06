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

type ResumeLogRepo struct {
}

func (s *ResumeLogRepo) Create(ctx context.Context, resumeLog *domain.ResumeLog) error {
	return global.DB.WithContext(ctx).Create(&resumeLog).Error
}

func (s *ResumeLogRepo) Delete(ctx context.Context, resumeLog domain.ResumeLog) error {
	return global.DB.WithContext(ctx).Delete(&resumeLog).Error
}

func (s *ResumeLogRepo) DeleteByIds(ctx context.Context, ids request.Ids) error {
	return global.DB.WithContext(ctx).Delete(&[]domain.ResumeLog{}, ids.Ids).Error
}

func (s *ResumeLogRepo) Update(ctx context.Context, resumeLog domain.ResumeLog) error {
	if resumeLog.ID == 0 {
		return errors.New(fmt.Sprintf("missing %s.id", "resumeLog"))
	}
	return global.DB.WithContext(ctx).Model(&domain.ResumeLog{Model: model.Model{ID: resumeLog.ID}}).Updates(&resumeLog).Error
}

func (s *ResumeLogRepo) Find(ctx context.Context, resumeLog domain.ResumeLog) (domain.ResumeLog, error) {
	db := global.DB.WithContext(ctx).Model(&domain.ResumeLog{})
	// TODO：条件过滤

	res := db.First(&resumeLog)

	return resumeLog, res.Error
}

func (s *ResumeLogRepo) List(ctx context.Context, page domain.PageResumeLogSearch) ([]domain.ResumeLog, int64, error) {
	var (
		resumeLogList []domain.ResumeLog
		count         int64
		err           error
	)
	// db
	db := global.DB.WithContext(ctx).Model(&domain.ResumeLog{})

	// TODO：条件过滤

	err = db.Count(&count).Scopes(scopes.Paginate(page.PageSearch), scopes.OrderBy(page.OrderBy)).Find(&resumeLogList).Error

	return resumeLogList, count, err
}
