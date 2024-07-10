package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools"
	"github.com/Madou-Shinni/go-logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrorUserExist = errors.New("account already exist")
	ErrorAccount   = errors.New("账号或密码错误")
)

var sysRoleService = NewSysRoleService()

// 定义接口
type SysUserRepo interface {
	Create(ctx context.Context, sysUser domain.SysUser) error
	Delete(ctx context.Context, sysUser domain.SysUser) error
	Update(ctx context.Context, sysUser map[string]interface{}) error
	Find(ctx context.Context, sysUser domain.SysUser) (domain.SysUser, error)
	List(ctx context.Context, page domain.PageSysUserSearch) ([]domain.SysUser, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type SysUserService struct {
	repo SysUserRepo
}

func NewSysUserService() *SysUserService {
	return &SysUserService{repo: &data.SysUserRepo{}}
}

func (s *SysUserService) Add(ctx context.Context, sysUser domain.SysUser) error {
	// 3.持久化入库
	db := global.DB.Model(&domain.SysUser{})
	err := db.Where("account = ?", sysUser.Account).First(&domain.SysUser{}).Error
	if err == nil {
		return ErrorUserExist
	}

	if err := s.repo.Create(ctx, sysUser); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(sysUser)", zap.Error(err), zap.Any("domain.SysUser", sysUser))
		return err
	}

	return nil
}

func (s *SysUserService) Delete(ctx context.Context, sysUser domain.SysUser) error {
	if err := s.repo.Delete(ctx, sysUser); err != nil {
		logger.Error("s.repo.Delete(sysUser)", zap.Error(err), zap.Any("domain.SysUser", sysUser))
		return err
	}

	return nil
}

func (s *SysUserService) Update(ctx context.Context, sysUser map[string]interface{}) error {
	if err := s.repo.Update(ctx, sysUser); err != nil {
		logger.Error("s.repo.Update(sysUser)", zap.Error(err), zap.Any("domain.SysUser", sysUser))
		return err
	}

	return nil
}

func (s *SysUserService) Find(ctx context.Context, sysUser domain.SysUser) (domain.SysUser, error) {
	res, err := s.repo.Find(ctx, sysUser)

	if err != nil {
		logger.Error("s.repo.Find(sysUser)", zap.Error(err), zap.Any("domain.SysUser", sysUser))
		return res, err
	}

	for i, v := range res.Roles {
		tree, err := sysRoleService.GetRoleTree(global.DB, v.ID)
		if err != nil {
			return res, err
		}
		res.Roles[i].Children = tree
	}

	return res, nil
}

func (s *SysUserService) List(ctx context.Context, page domain.PageSysUserSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageSysUserSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *SysUserService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}

func (s *SysUserService) Login(ctx context.Context, user domain.LoginReq) (interface{}, error) {
	// 查询用户
	var sysUser domain.SysUser
	err := global.DB.WithContext(ctx).Model(&domain.SysUser{}).First(&sysUser, "account = ?", user.Account).Error
	if err != nil {
		logger.Error("s.Login()", zap.Error(err), zap.Any("domain.SysUser", user))
		return nil, ErrorAccount
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(sysUser.Password), []byte(user.Password))
	if err != nil {
		return nil, ErrorAccount
	}

	// 生成token
	mp := jwt.MapClaims{
		tools.UserIdKey: sysUser.ID,
		tools.RoleIdKey: sysUser.DefaultRole,
		tools.ExpKey:    time.Duration(conf.Conf.JwtConfig.AccessExpire) * time.Second,
	}
	token, err := tools.GenToken(mp, conf.Conf.JwtConfig.Secret)
	if err != nil {
		logger.Error("token生成失败", zap.Error(err), zap.Any("jwt.MapClaims", mp))
		return nil, fmt.Errorf(constant.CODE_ERR_BUSY.Msg())
	}

	return token, nil
}
