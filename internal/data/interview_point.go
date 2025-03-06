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

type InterviewPointRepo struct {
}

func (s *InterviewPointRepo) Create(ctx context.Context, interviewPoint *domain.InterviewPoint) error {
	return global.DB.WithContext(ctx).Create(&interviewPoint).Error
}

func (s *InterviewPointRepo) Delete(ctx context.Context, interviewPoint domain.InterviewPoint) error {
	return global.DB.WithContext(ctx).Delete(&interviewPoint).Error
}

func (s *InterviewPointRepo) DeleteByIds(ctx context.Context, ids request.Ids) error {
	return global.DB.WithContext(ctx).Delete(&[]domain.InterviewPoint{}, ids.Ids).Error
}

func (s *InterviewPointRepo) Update(ctx context.Context, interviewPoint domain.InterviewPoint) error {
	if interviewPoint.ID == 0 {
		return errors.New(fmt.Sprintf("missing %s.id", "interviewPoint"))
	}
	return global.DB.WithContext(ctx).Model(&domain.InterviewPoint{Model: model.Model{ID: interviewPoint.ID}}).Updates(&interviewPoint).Error
}

func (s *InterviewPointRepo) Find(ctx context.Context, interviewPoint domain.InterviewPoint) (domain.InterviewPoint, error) {
	db := global.DB.WithContext(ctx).Model(&domain.InterviewPoint{})
	// TODO：条件过滤

	res := db.First(&interviewPoint)

	return interviewPoint, res.Error
}

func (s *InterviewPointRepo) List(ctx context.Context, page domain.PageInterviewPointSearch) ([]domain.InterviewPoint, int64, error) {
	var (
		interviewPointList []domain.InterviewPoint
		count              int64
		err                error
	)
	// db
	db := global.DB.WithContext(ctx).Model(&domain.InterviewPoint{})

	// TODO：条件过滤

	err = db.Count(&count).Scopes(scopes.Paginate(page.PageSearch), scopes.OrderBy(page.OrderBy)).Find(&interviewPointList).Error

	return interviewPointList, count, err
}
