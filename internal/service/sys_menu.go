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
type SysMenuRepo interface {
	Create(ctx context.Context, sysMenu domain.SysMenu) error
	Delete(ctx context.Context, sysMenu domain.SysMenu) error
	Update(ctx context.Context, sysMenu map[string]interface{}) error
	Find(ctx context.Context, sysMenu domain.SysMenu) (domain.SysMenu, error)
	List(ctx context.Context, page domain.PageSysMenuSearch) ([]domain.SysMenu, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type SysMenuService struct {
	repo SysMenuRepo
}

func NewSysMenuService() *SysMenuService {
	return &SysMenuService{repo: &data.SysMenuRepo{}}
}

func (s *SysMenuService) Add(ctx context.Context, sysMenu domain.SysMenu) error {
	// 3.持久化入库
	if err := s.repo.Create(ctx, sysMenu); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(sysMenu)", zap.Error(err), zap.Any("domain.SysMenu", sysMenu))
		return err
	}

	return nil
}

func (s *SysMenuService) Delete(ctx context.Context, sysMenu domain.SysMenu) error {
	if err := s.repo.Delete(ctx, sysMenu); err != nil {
		logger.Error("s.repo.Delete(sysMenu)", zap.Error(err), zap.Any("domain.SysMenu", sysMenu))
		return err
	}

	return nil
}

func (s *SysMenuService) Update(ctx context.Context, sysMenu map[string]interface{}) error {
	if err := s.repo.Update(ctx, sysMenu); err != nil {
		logger.Error("s.repo.Update(sysMenu)", zap.Error(err), zap.Any("domain.SysMenu", sysMenu))
		return err
	}

	return nil
}

func (s *SysMenuService) Find(ctx context.Context, sysMenu domain.SysMenu) (domain.SysMenu, error) {
	res, err := s.repo.Find(ctx, sysMenu)

	if err != nil {
		logger.Error("s.repo.Find(sysMenu)", zap.Error(err), zap.Any("domain.SysMenu", sysMenu))
		return res, err
	}

	return res, nil
}

func (s *SysMenuService) List(ctx context.Context, page domain.PageSysMenuSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageSysMenuSearch", page))
		return pageRes, err
	}

	for i, v := range data {
		tree, err := s.GetMenuTree(global.DB.WithContext(ctx), v.ID)
		if err != nil {
			return pageRes, err
		}
		data[i].Children = tree
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *SysMenuService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}

func (s *SysMenuService) RoleList(ctx context.Context, rid uint) ([]domain.SysMenu, error) {
	var list []domain.SysMenu

	err := global.DB.WithContext(ctx).Model(&domain.SysRole{Model: model.Model{ID: rid}}).
		Where("parent_id = ?", 0).
		Association("Menus").
		Find(&list)
	if err != nil {
		return nil, err
	}

	// 构建树形菜单
	result := buildTreeMenu(list, 0)

	return result, nil
}

// 转换为树形菜单
func buildTreeMenu(menus []domain.SysMenu, parentID uint) []domain.SysMenu {
	var tree []domain.SysMenu
	menuMap := make(map[uint][]domain.SysMenu)

	// 按 parent_id 分组
	for _, v := range menus {
		menuMap[v.ParentID] = append(menuMap[v.ParentID], v)
	}

	// 递归构建树
	var build func(parentID uint) []domain.SysMenu
	build = func(parentID uint) []domain.SysMenu {
		children := menuMap[parentID]
		for i, child := range children {
			children[i].Children = build(child.ID) // 递归查找子节点
		}
		return children
	}

	tree = build(parentID)
	return tree
}

func (s *SysMenuService) SetRoleList(ctx context.Context, sysRole domain.SysRole) error {
	var list = sysRole.Menus

	err := global.DB.WithContext(ctx).Model(&domain.SysRole{Model: model.Model{ID: sysRole.ID}}).
		Association("Menus").
		Replace(list)
	if err != nil {
		return err
	}

	return err
}

// GetMenuTree 递归树
func (s *SysMenuService) GetMenuTree(db *gorm.DB, parentID uint) ([]domain.SysMenu, error) {
	var menus []domain.SysMenu
	if err := db.Where("parent_id = ?", parentID).Find(&menus).Error; err != nil {
		return nil, err
	}

	for i := range menus {
		children, err := s.GetMenuTree(db, menus[i].ID)
		if err != nil {
			return nil, err
		}
		menus[i].Children = children
	}

	return menus, nil
}
