package service

import (
	"context"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
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

	for i, v := range data {
		tree, err := s.GetRoleTree(global.DB, v.ID)
		if err != nil {
			return pageRes, err
		}
		data[i].Children = tree
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

func (s *SysRoleService) SetUserRoleList(ctx context.Context, sysUser domain.SysUser) error {
	var sysRoles = sysUser.Roles
	var ok bool
	// 查询用户角色列表
	err := global.DB.WithContext(ctx).Model(&domain.SysUser{}).First(&sysUser, "id = ?", sysUser.ID).Error
	if err != nil {
		return err
	}

	for _, v := range sysRoles {
		if v.ID == sysUser.DefaultRole {
			ok = true
		}
	}

	err = global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if !ok {
			// 设置角色不包含 用户目前的默认角色
			// 用户默认角色为空
			var defaultRole uint
			if len(sysRoles) > 0 {
				defaultRole = sysRoles[0].ID
			}
			err = tx.Model(&domain.SysUser{}).Where("id = ?", sysUser.ID).UpdateColumn("default_role", defaultRole).Error
			if err != nil {
				return err
			}
		}

		err = tx.Model(&domain.SysUser{Model: model.Model{ID: sysUser.ID}}).
			Association("Roles").
			Replace(sysUser.Roles)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return err
}

// GetRoleTree 递归树
func (s *SysRoleService) GetRoleTree(db *gorm.DB, parentID uint) ([]domain.SysRole, error) {
	var menus []domain.SysRole
	if err := db.Where("parent_id = ?", parentID).Find(&menus).Error; err != nil {
		return nil, err
	}

	for i := range menus {
		children, err := s.GetRoleTree(db, menus[i].ID)
		if err != nil {
			return nil, err
		}
		menus[i].Children = children
	}

	return menus, nil
}
