package service

import (
	"context"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/excel"
	"github.com/Madou-Shinni/go-logger"
	"go.uber.org/zap"
)

type DataImportTemplateReq struct {
	Category string `json:"category" form:"category"` // 类型
}

type DemoExcelTpl struct {
	Name     string           `excel:"姓名"`
	Age      int              `excel:"年龄"`
	BirthDay *model.LocalTime `excel:"生日"`
}

// 定义接口
type DataImportRepo interface {
	Create(ctx context.Context, dataImport *domain.DataImport) error
	Delete(ctx context.Context, dataImport domain.DataImport) error
	Update(ctx context.Context, dataImport domain.DataImport) error
	Find(ctx context.Context, dataImport domain.DataImport) (domain.DataImport, error)
	List(ctx context.Context, page domain.PageDataImportSearch) ([]domain.DataImport, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type DataImportService struct {
	repo DataImportRepo
}

func NewDataImportService() *DataImportService {
	return &DataImportService{repo: &data.DataImportRepo{}}
}

func (s *DataImportService) Add(ctx context.Context, dataImport domain.DataImport) error {
	// 3.持久化入库
	if err := s.repo.Create(ctx, &dataImport); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(dataImport)", zap.Error(err), zap.Any("domain.DataImport", dataImport))
		return err
	}

	return nil
}

func (s *DataImportService) Delete(ctx context.Context, dataImport domain.DataImport) error {
	if err := s.repo.Delete(ctx, dataImport); err != nil {
		logger.Error("s.repo.Delete(dataImport)", zap.Error(err), zap.Any("domain.DataImport", dataImport))
		return err
	}

	return nil
}

func (s *DataImportService) Update(ctx context.Context, dataImport domain.DataImport) error {
	if err := s.repo.Update(ctx, dataImport); err != nil {
		logger.Error("s.repo.Update(dataImport)", zap.Error(err), zap.Any("domain.DataImport", dataImport))
		return err
	}

	return nil
}

func (s *DataImportService) Find(ctx context.Context, dataImport domain.DataImport) (domain.DataImport, error) {
	res, err := s.repo.Find(ctx, dataImport)

	if err != nil {
		logger.Error("s.repo.Find(dataImport)", zap.Error(err), zap.Any("domain.DataImport", dataImport))
		return res, err
	}

	return res, nil
}

func (s *DataImportService) List(ctx context.Context, page domain.PageDataImportSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageDataImportSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *DataImportService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}

func (s *DataImportService) Template(ctx context.Context, req DataImportTemplateReq) ([]byte, error) {
	var err error
	tool := excel.NewExcelTool("Sheet1")
	switch req.Category {
	case constants.DataImportCategoryDemo:
		err = demoExcelTpl(ctx, tool)
	default:
		return nil, fmt.Errorf("暂不支持该类型")
	}

	if err != nil {
		return nil, err
	}

	buffer, err := tool.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	//tool.SaveAs(fmt.Sprintf("%s_%s.xlsx", req.Category, uuid.NewString()))

	return buffer.Bytes(), nil
}

func (s *DataImportService) Import(ctx context.Context, req domain.DataImport) (interface{}, error) {
	if err := global.DB.WithContext(ctx).Create(&req).Error; err != nil {
		return nil, err
	}
	err := global.Producer.NewTask(constants.QueueDataImport, req)
	if err != nil {
		global.DB.WithContext(ctx).Model(&req).Update("status", constants.DataImportStatusFailed)
		return nil, err
	}
	return nil, nil
}

func demoExcelTpl(ctx context.Context, tool *excel.ExcelTool, fns ...func(ctx context.Context, tool *excel.ExcelTool) error) error {
	tool.Model(&DemoExcelTpl{})

	for _, f := range fns {
		if err := f(ctx, tool); err != nil {
			return err
		}
	}

	return tool.Flush()
}
