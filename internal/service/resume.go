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
type ResumeRepo interface {
	Create(ctx context.Context, resume *domain.Resume) error
	Delete(ctx context.Context, resume domain.Resume) error
	Update(ctx context.Context, resume domain.Resume) error
	Find(ctx context.Context, resume domain.Resume) (domain.Resume, error)
	List(ctx context.Context, page domain.PageResumeSearch) ([]domain.Resume, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type ResumeService struct {
	repo ResumeRepo
}

func NewResumeService() *ResumeService {
	return &ResumeService{repo: &data.ResumeRepo{}}
}

func (s *ResumeService) Add(ctx context.Context, resume domain.Resume) error {
	// 3.持久化入库
	if err := s.repo.Create(ctx, &resume); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(resume)", zap.Error(err), zap.Any("domain.Resume", resume))
		return err
	}

	return nil
}

func (s *ResumeService) Delete(ctx context.Context, resume domain.Resume) error {
	if err := s.repo.Delete(ctx, resume); err != nil {
		logger.Error("s.repo.Delete(resume)", zap.Error(err), zap.Any("domain.Resume", resume))
		return err
	}

	return nil
}

func (s *ResumeService) Update(ctx context.Context, resume domain.Resume) error {
	if err := s.repo.Update(ctx, resume); err != nil {
		logger.Error("s.repo.Update(resume)", zap.Error(err), zap.Any("domain.Resume", resume))
		return err
	}

	return nil
}

func (s *ResumeService) Find(ctx context.Context, resume domain.Resume) (domain.Resume, error) {
	res, err := s.repo.Find(ctx, resume)

	if err != nil {
		logger.Error("s.repo.Find(resume)", zap.Error(err), zap.Any("domain.Resume", resume))
		return res, err
	}

	return res, nil
}

func (s *ResumeService) List(ctx context.Context, page domain.PageResumeSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageResumeSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *ResumeService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}
