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
type DemoRepo interface {
	Create(ctx context.Context, demo domain.Demo) error
	Delete(ctx context.Context, demo domain.Demo) error
	Update(ctx context.Context, demo domain.Demo) error
	Find(ctx context.Context, demo domain.Demo) (domain.Demo, error)
	List(ctx context.Context, page domain.PageDemoSearch) ([]domain.Demo, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type DemoService struct {
	repo DemoRepo
}

func NewDemoService() *DemoService {
	return &DemoService{repo: &data.DemoRepo{}}
}

func (s *DemoService) Add(ctx context.Context, demo domain.Demo) error {
	// 3.持久化入库
	if err := s.repo.Create(ctx, demo); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Delete(ctx context.Context, demo domain.Demo) error {
	if err := s.repo.Delete(ctx, demo); err != nil {
		logger.Error("s.repo.Delete(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Update(ctx context.Context, demo domain.Demo) error {
	if err := s.repo.Update(ctx, demo); err != nil {
		logger.Error("s.repo.Update(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Find(ctx context.Context, demo domain.Demo) (domain.Demo, error) {
	res, err := s.repo.Find(ctx, demo)

	if err != nil {
		logger.Error("s.repo.Find(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return res, err
	}

	return res, nil
}

func (s *DemoService) List(ctx context.Context, page domain.PageDemoSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageDemoSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *DemoService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}
