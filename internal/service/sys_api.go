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
type SysApiRepo interface {
	Create(ctx context.Context, sysApi domain.SysApi) error
	Delete(ctx context.Context, sysApi domain.SysApi) error
	Update(ctx context.Context, sysApi map[string]interface{}) error
	Find(ctx context.Context, sysApi domain.SysApi) (domain.SysApi, error)
	List(ctx context.Context, page domain.PageSysApiSearch) ([]domain.SysApi, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type SysApiService struct {
	repo SysApiRepo
}

func NewSysApiService() *SysApiService {
	return &SysApiService{repo: &data.SysApiRepo{}}
}

func (s *SysApiService) Add(ctx context.Context, sysApi domain.SysApi) error {
	// 3.持久化入库
	if err := s.repo.Create(ctx, sysApi); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(sysApi)", zap.Error(err), zap.Any("domain.SysApi", sysApi))
		return err
	}

	return nil
}

func (s *SysApiService) Delete(ctx context.Context, sysApi domain.SysApi) error {
	if err := s.repo.Delete(ctx, sysApi); err != nil {
		logger.Error("s.repo.Delete(sysApi)", zap.Error(err), zap.Any("domain.SysApi", sysApi))
		return err
	}

	return nil
}

func (s *SysApiService) Update(ctx context.Context, sysApi map[string]interface{}) error {
	if err := s.repo.Update(ctx, sysApi); err != nil {
		logger.Error("s.repo.Update(sysApi)", zap.Error(err), zap.Any("domain.SysApi", sysApi))
		return err
	}

	return nil
}

func (s *SysApiService) Find(ctx context.Context, sysApi domain.SysApi) (domain.SysApi, error) {
	res, err := s.repo.Find(ctx, sysApi)

	if err != nil {
		logger.Error("s.repo.Find(sysApi)", zap.Error(err), zap.Any("domain.SysApi", sysApi))
		return res, err
	}

	return res, nil
}

func (s *SysApiService) List(ctx context.Context, page domain.PageSysApiSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageSysApiSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *SysApiService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}
