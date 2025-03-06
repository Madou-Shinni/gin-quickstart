package service

import (
	"context"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/go-logger"
	"go.uber.org/zap"
)

// 定义接口
type InterviewPointRepo interface {
	Create(ctx context.Context, interviewPoint *domain.InterviewPoint) error
	Delete(ctx context.Context, interviewPoint domain.InterviewPoint) error
	Update(ctx context.Context, interviewPoint domain.InterviewPoint) error
	Find(ctx context.Context, interviewPoint domain.InterviewPoint) (domain.InterviewPoint, error)
	List(ctx context.Context, page domain.PageInterviewPointSearch) ([]domain.InterviewPoint, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type InterviewPointService struct {
	repo InterviewPointRepo
}

func NewInterviewPointService() *InterviewPointService {
	return &InterviewPointService{repo: &data.InterviewPointRepo{}}
}

func (s *InterviewPointService) Add(ctx context.Context, interviewPoint domain.InterviewPoint) error {
	// 3.持久化入库
	if err := s.repo.Create(ctx, &interviewPoint); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(interviewPoint)", zap.Error(err), zap.Any("domain.InterviewPoint", interviewPoint))
		return err
	}

	return nil
}

func (s *InterviewPointService) Delete(ctx context.Context, interviewPoint domain.InterviewPoint) error {
	if err := s.repo.Delete(ctx, interviewPoint); err != nil {
		logger.Error("s.repo.Delete(interviewPoint)", zap.Error(err), zap.Any("domain.InterviewPoint", interviewPoint))
		return err
	}

	return nil
}

func (s *InterviewPointService) Update(ctx context.Context, interviewPoint domain.InterviewPoint) error {
	if err := s.repo.Update(ctx, interviewPoint); err != nil {
		logger.Error("s.repo.Update(interviewPoint)", zap.Error(err), zap.Any("domain.InterviewPoint", interviewPoint))
		return err
	}

	return nil
}

func (s *InterviewPointService) Find(ctx context.Context, interviewPoint domain.InterviewPoint) (domain.InterviewPoint, error) {
	res, err := s.repo.Find(ctx, interviewPoint)

	if err != nil {
		logger.Error("s.repo.Find(interviewPoint)", zap.Error(err), zap.Any("domain.InterviewPoint", interviewPoint))
		return res, err
	}

	return res, nil
}

func (s *InterviewPointService) List(ctx context.Context, page domain.PageInterviewPointSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageInterviewPointSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *InterviewPointService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}
