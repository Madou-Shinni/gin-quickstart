package service

import (
	"context"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/go-logger"
	"github.com/casbin/casbin/v2"
	csmodel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"sync"
)

const (
	rbac_models = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch4(r.obj, p.obj) && regexMatch(r.act, p.act)
`
)

var (
	once sync.Once
	e    *casbin.Enforcer
)

// 定义接口
type SysCasbinRepo interface {
	Create(ctx context.Context, sysCasbin domain.SysCasbin) error
	Delete(ctx context.Context, sysCasbin domain.SysCasbin) error
	Update(ctx context.Context, sysCasbin map[string]interface{}) error
	Find(ctx context.Context, sysCasbin domain.SysCasbin) (domain.SysCasbin, error)
	List(ctx context.Context, page domain.PageSysCasbinSearch) ([]domain.SysCasbin, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type SysCasbinService struct {
	repo SysCasbinRepo
	e    func() *casbin.Enforcer
}

func NewSysCasbinService() *SysCasbinService {
	s := &SysCasbinService{repo: &data.SysCasbinRepo{}}
	s.e = Casbin
	return s
}

func (s *SysCasbinService) AddUserRoles(ctx context.Context, req domain.UserRolesReq) error {
	ucs := constant.GetCasbinUserKey(req.UserID)
	// 删除
	_, err := s.e().DeleteRolesForUser(ucs)
	if err != nil {
		return err
	}
	// 添加用户角色
	slice := make([]string, 0, len(req.Roles))
	for _, role := range req.Roles {
		rcs := constant.GetCasbinRoleKey(role)
		slice = append(slice, rcs)
	}

	_, err = s.e().AddRolesForUser(ucs, slice)
	if err != nil {
		return err
	}

	return nil
}

func (s *SysCasbinService) AddRoleRoles(ctx context.Context, req domain.RoleRolesReq) error {
	// 删除角色
	rcs := constant.GetCasbinRoleKey(req.Role)
	_, err := s.e().DeleteRolesForUser(rcs)
	if err != nil {
		return err
	}
	// 添加角色
	slice := make([]string, 0, len(req.Roles))
	for _, role := range req.Roles {
		rr := constant.GetCasbinRoleKey(role)
		slice = append(slice, rr)
	}

	_, err = s.e().AddRolesForUser(rcs, slice)
	if err != nil {
		return err
	}

	return nil
}

func (s *SysCasbinService) AddRolePermissions(ctx context.Context, req domain.RolePermissionsReq) error {
	rcs := constant.GetCasbinRoleKey(req.Role)
	// 删除权限
	_, err := s.e().DeletePermissionsForUser(rcs)
	if err != nil {
		return err
	}
	// 添加权限
	for k, v := range req.Permissions {
		_, err = s.e().AddPermissionsForUser(rcs, []string{k, v})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SysCasbinService) Delete(ctx context.Context, sysCasbin domain.SysCasbin) error {
	if err := s.repo.Delete(ctx, sysCasbin); err != nil {
		logger.Error("s.repo.Delete(sysCasbin)", zap.Error(err), zap.Any("domain.SysCasbin", sysCasbin))
		return err
	}

	return nil
}

func (s *SysCasbinService) Update(ctx context.Context, sysCasbin map[string]interface{}) error {
	if err := s.repo.Update(ctx, sysCasbin); err != nil {
		logger.Error("s.repo.Update(sysCasbin)", zap.Error(err), zap.Any("domain.SysCasbin", sysCasbin))
		return err
	}

	return nil
}

func (s *SysCasbinService) Find(ctx context.Context, sysCasbin domain.SysCasbin) (domain.SysCasbin, error) {
	res, err := s.repo.Find(ctx, sysCasbin)

	if err != nil {
		logger.Error("s.repo.Find(sysCasbin)", zap.Error(err), zap.Any("domain.SysCasbin", sysCasbin))
		return res, err
	}

	return res, nil
}

func (s *SysCasbinService) List(ctx context.Context, page domain.PageSysCasbinSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageSysCasbinSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *SysCasbinService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}

func Casbin() *casbin.Enforcer {
	once.Do(func() {
		m, err := csmodel.NewModelFromString(rbac_models)
		if err != nil {
			logger.Error("Casbin NewModelFromString", zap.Error(err))
			return
		}

		a, err := gormadapter.NewAdapterByDB(global.DB)
		if err != nil {
			logger.Error("Casbin NewAdapterByDB", zap.Error(err))
			return
		}
		en, err := casbin.NewEnforcer(m, a)
		if err != nil {
			logger.Error("Casbin NewModelFromString", zap.Error(err))
			return
		}

		e = en
	})

	return e
}
