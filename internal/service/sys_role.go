package service

import (
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
	Create(sysRole domain.SysRole) error
	Delete(sysRole domain.SysRole) error
	Update(sysRole map[string]interface{}) error
	Find(sysRole domain.SysRole) (domain.SysRole, error)
	List(page domain.PageSysRoleSearch) ([]domain.SysRole, int64, error)
	DeleteByIds(ids request.Ids) error
}

type SysRoleService struct {
	repo          SysRoleRepo
	casbinService *SysCasbinService
}

func NewSysRoleService() *SysRoleService {
	return &SysRoleService{repo: &data.SysRoleRepo{}, casbinService: NewSysCasbinService()}
}

func (s *SysRoleService) Add(sysRole domain.SysRole) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&domain.SysRole{}).Create(&sysRole).Error
		if err != nil {
			return err
		}

		if sysRole.ParentID != 0 {
			// 子角色
			err = s.casbinService.AddRoleRoles(domain.RoleRolesReq{
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

func (s *SysRoleService) Delete(sysRole domain.SysRole) error {
	if err := s.repo.Delete(sysRole); err != nil {
		logger.Error("s.repo.Delete(sysRole)", zap.Error(err), zap.Any("domain.SysRole", sysRole))
		return err
	}

	return nil
}

func (s *SysRoleService) Update(sysRole map[string]interface{}) error {
	if err := s.repo.Update(sysRole); err != nil {
		logger.Error("s.repo.Update(sysRole)", zap.Error(err), zap.Any("domain.SysRole", sysRole))
		return err
	}

	return nil
}

func (s *SysRoleService) Find(sysRole domain.SysRole) (domain.SysRole, error) {
	res, err := s.repo.Find(sysRole)

	if err != nil {
		logger.Error("s.repo.Find(sysRole)", zap.Error(err), zap.Any("domain.SysRole", sysRole))
		return res, err
	}

	return res, nil
}

func (s *SysRoleService) List(page domain.PageSysRoleSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageSysRoleSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *SysRoleService) DeleteByIds(ids request.Ids) error {
	if err := s.repo.DeleteByIds(ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}
