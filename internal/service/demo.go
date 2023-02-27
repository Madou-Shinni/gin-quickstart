package service

import (
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/letter"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/snowflake"
	"github.com/Madou-Shinni/go-logger"
	"go.uber.org/zap"
)

// 定义接口
type DemoRepo interface {
	Create(demo domain.Demo) error
	Delete(demo domain.Demo) error
	Update(demo domain.Demo) error
	Find(demo domain.Demo) (domain.Demo, error)
	List(page domain.PageDemoSearch) ([]domain.Demo, error)
	Count() (int64, error)
	DeleteByIds(ids request.Ids) error
}

type DemoService struct {
	repo DemoRepo
}

func NewDemoService() *DemoService {
	return &DemoService{repo: data.NewDemoRepo()}
}

func (s *DemoService) Add(demo domain.Demo) error {
	// 1.生成唯一标识
	// 因为我们在全局初始化的时候已经初始化了雪花算法的机器节点
	// 所以我们可以直接使用
	did := snowflake.GenerateID()

	// 2.生成code 20位
	code := letter.GenerateCode(20)

	demo.Did = did
	demo.Code = code

	// 3.持久化入库
	if err := s.repo.Create(demo); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Delete(demo domain.Demo) error {
	if err := s.repo.Delete(demo); err != nil {
		logger.Error("s.repo.Delete(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Update(demo domain.Demo) error {
	if err := s.repo.Update(demo); err != nil {
		logger.Error("s.repo.Update(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Find(demo domain.Demo) (domain.Demo, error) {
	res, err := s.repo.Find(demo)

	if err != nil {
		logger.Error("s.repo.Find(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return res, err
	}

	return res, nil
}

func (s *DemoService) List(page domain.PageDemoSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, err := s.repo.List(page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageDemoSearch", page))
		return pageRes, err
	}

	count, err := s.repo.Count()
	if err != nil {
		logger.Error("s.repo.Count()", zap.Error(err))
		return pageRes, err
	}

	pageRes.Data = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *DemoService) DeleteByIds(ids request.Ids) error {
	if err := s.repo.DeleteByIds(ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}
