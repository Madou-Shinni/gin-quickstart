package service

import (
	"context"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/go-logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 定义接口
type SysRoleRepo interface {
	Create(ctx context.Context, sysRole domain.SysRole) error
	Delete(ctx context.Context, sysRole domain.SysRole) error
	Update(ctx context.Context, sysRole map[string]interface{}) error
	Find(ctx context.Context, sysRole domain.SysRole) (domain.SysRole, error)
	List(ctx context.Context, page domain.PageSysRoleSearch) ([]domain.SysRole, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type SysRoleService struct {
	repo          SysRoleRepo
	casbinService *SysCasbinService
}

func NewSysRoleService() *SysRoleService {
	return &SysRoleService{repo: &data.SysRoleRepo{}, casbinService: NewSysCasbinService()}
}

func (s *SysRoleService) Add(ctx context.Context, sysRole domain.SysRole) error {
	err := global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&domain.SysRole{}).Create(&sysRole).Error
		if err != nil {
			return err
		}

		if sysRole.ParentID != 0 {
			// 子角色
			err = s.casbinService.AddRoleRoles(ctx, domain.RoleRolesReq{
				Role:  sysRole.ParentID,
				Roles: []uint{sysRole.ID},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SysRoleService) Delete(ctx context.Context, sysRole domain.SysRole) error {
	if err := s.repo.Delete(ctx, sysRole); err != nil {
		logger.Error("s.repo.Delete(sysRole)", zap.Error(err), zap.Any("domain.SysRole", sysRole))
		return err
	}

	return nil
}

func (s *SysRoleService) Update(ctx context.Context, sysRole map[string]interface{}) error {
	if err := s.repo.Update(ctx, sysRole); err != nil {
		logger.Error("s.repo.Update(sysRole)", zap.Error(err), zap.Any("domain.SysRole", sysRole))
		return err
	}

	return nil
}

func (s *SysRoleService) Find(ctx context.Context, sysRole domain.SysRole) (domain.SysRole, error) {
	res, err := s.repo.Find(ctx, sysRole)

	if err != nil {
		logger.Error("s.repo.Find(sysRole)", zap.Error(err), zap.Any("domain.SysRole", sysRole))
		return res, err
	}

	return res, nil
}

func (s *SysRoleService) List(ctx context.Context, page domain.PageSysRoleSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageSysRoleSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *SysRoleService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}
