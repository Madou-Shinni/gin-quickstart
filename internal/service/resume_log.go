package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/distance"
	"github.com/Madou-Shinni/go-logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ResumePushReq struct {
	RID         uint    `json:"rid"`          // 简历id
	SelectedIID *uint   `json:"selected_iid"` // 用户选择的面试点
	Latitude    float64 `json:"latitude"`     // 纬度
	Longitude   float64 `json:"longitude"`    // 经度
	NearLimit   float64 `json:"near_limit"`   // 附近范围 km
}

// 定义接口
type ResumeLogRepo interface {
	Create(ctx context.Context, resumeLog *domain.ResumeLog) error
	Delete(ctx context.Context, resumeLog domain.ResumeLog) error
	Update(ctx context.Context, resumeLog domain.ResumeLog) error
	Find(ctx context.Context, resumeLog domain.ResumeLog) (domain.ResumeLog, error)
	List(ctx context.Context, page domain.PageResumeLogSearch) ([]domain.ResumeLog, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type ResumeLogService struct {
	repo ResumeLogRepo
}

func NewResumeLogService() *ResumeLogService {
	return &ResumeLogService{repo: &data.ResumeLogRepo{}}
}

func (s *ResumeLogService) Add(ctx context.Context, resumeLog domain.ResumeLog) error {
	// 3.持久化入库
	if err := s.repo.Create(ctx, &resumeLog); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(resumeLog)", zap.Error(err), zap.Any("domain.ResumeLog", resumeLog))
		return err
	}

	return nil
}

func (s *ResumeLogService) Delete(ctx context.Context, resumeLog domain.ResumeLog) error {
	if err := s.repo.Delete(ctx, resumeLog); err != nil {
		logger.Error("s.repo.Delete(resumeLog)", zap.Error(err), zap.Any("domain.ResumeLog", resumeLog))
		return err
	}

	return nil
}

func (s *ResumeLogService) Update(ctx context.Context, resumeLog domain.ResumeLog) error {
	if err := s.repo.Update(ctx, resumeLog); err != nil {
		logger.Error("s.repo.Update(resumeLog)", zap.Error(err), zap.Any("domain.ResumeLog", resumeLog))
		return err
	}

	return nil
}

func (s *ResumeLogService) Find(ctx context.Context, resumeLog domain.ResumeLog) (domain.ResumeLog, error) {
	res, err := s.repo.Find(ctx, resumeLog)

	if err != nil {
		logger.Error("s.repo.Find(resumeLog)", zap.Error(err), zap.Any("domain.ResumeLog", resumeLog))
		return res, err
	}

	return res, nil
}

func (s *ResumeLogService) List(ctx context.Context, page domain.PageResumeLogSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageResumeLogSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *ResumeLogService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}

func (s *ResumeLogService) Push(ctx context.Context, req ResumePushReq) (interface{}, error) {
	// 匹配面试点
	ip, err := s.MatchIp(ctx, req)
	if err != nil {
		return nil, err
	}

	err = global.DB.Tx(ctx, func(ctx context.Context) error {
		// 添加推送记录
		if err = s.Add(ctx, domain.ResumeLog{
			RID: req.RID,
			IID: ip.ID,
		}); err != nil {
			return err
		}

		// 面试点计数
		return global.DB.WithContext(ctx).Model(&domain.InterviewPoint{Model: model.Model{ID: ip.ID}}).
			UpdateColumn("rider_count", gorm.Expr("rider_count + 1")).Error
	})
	if err != nil {
		return nil, err
	}

	return ip, nil
}

func (s *ResumeLogService) MatchIp(ctx context.Context, req ResumePushReq) (*domain.InterviewPoint, error) {
	// 规则1：直接选择站点
	if req.SelectedIID != nil {
		if ip, err := validateIp(ctx, *req.SelectedIID); err != nil {
			return nil, err
		} else {
			return ip, nil
		}
	}

	// 规则2：区域匹配
	if req.Latitude != 0 && req.Longitude != 0 {
		if req.NearLimit == 0 {
			// 默认附近20公里
			req.NearLimit = 20
		}
		if ip, err := locateRegion(ctx, req.Latitude, req.Longitude, req.NearLimit); err != nil {
			return nil, err
		} else {
			return ip, nil
		}
	}

	// 规则3：均衡分配
	return balancedAssignment(ctx)
}

func validateIp(ctx context.Context, iid uint) (*domain.InterviewPoint, error) {
	var ip *domain.InterviewPoint
	if err := global.DB.WithContext(ctx).First(&ip, iid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("面点不存在")
		}
		return nil, err
	}
	return ip, nil
}

func locateRegion(ctx context.Context, lat, lng float64, nearLimit float64) (*domain.InterviewPoint, error) {
	var ipList []*domain.InterviewPoint
	minLat, maxLat, minLon, maxLon := distance.GetNearbyBoundingBox(lat, lng, nearLimit)
	db := global.DB.WithContext(ctx).Model(&domain.InterviewPoint{})
	// 计算查询距离
	db = db.Select("*", fmt.Sprintf(`6371 * ACOS(
        COS(RADIANS(%f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%f)) +
        SIN(RADIANS(%f)) * SIN(RADIANS(latitude))
    ) AS distance`, lat, lng, lat))
	// 粗略删选
	db = db.Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?", minLat, maxLat, minLon, maxLon)
	err := db.Having("distance <= ?", nearLimit). // 精确筛选（使用别名）xx km内 注意：这里Having条件放在Count查询之后 否则会报错
							Order("distance ASC").
							Order("rider_count ASC").
							Limit(1).
							Find(&ipList).Error
	if len(ipList) == 0 {
		return nil, errors.New("未找到匹配的面试点")
	}
	return ipList[0], err
}

func balancedAssignment(ctx context.Context) (*domain.InterviewPoint, error) {
	var ips []*domain.InterviewPoint
	if err := global.DB.WithContext(ctx).Order("rider_count ASC").Limit(1).Find(&ips).Error; err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, errors.New("未找到匹配的面试点")
	}

	return ips[0], nil
}
