package handle

import (
	"errors"
	"github.com/Madou-Shinni/gin-quickstart/common"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SysMenuHandle struct {
	s *service.SysMenuService
}

func NewSysMenuHandle() *SysMenuHandle {
	return &SysMenuHandle{s: service.NewSysMenuService()}
}

// Add 创建SysMenu
// @Tags     SysMenu
// @Summary  创建SysMenu
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysMenu true "创建SysMenu"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysMenu [post]
func (cl *SysMenuHandle) Add(c *gin.Context) {
	var sysMenu domain.SysMenu
	if err := c.ShouldBindJSON(&sysMenu); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), sysMenu); err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除SysMenu
// @Tags     SysMenu
// @Summary  删除SysMenu
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysMenu true "删除SysMenu"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysMenu [delete]
func (cl *SysMenuHandle) Delete(c *gin.Context) {
	var sysMenu domain.SysMenu
	if err := c.ShouldBindJSON(&sysMenu); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), sysMenu); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除SysMenu
// @Tags     SysMenu
// @Summary  批量删除SysMenu
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除SysMenu"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysMenu/delete-batch [delete]
func (cl *SysMenuHandle) DeleteByIds(c *gin.Context) {
	var ids request.Ids
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.DeleteByIds(c.Request.Context(), ids); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Update 修改SysMenu
// @Tags     SysMenu
// @Summary  修改SysMenu
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysMenu true "修改SysMenu"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysMenu [put]
func (cl *SysMenuHandle) Update(c *gin.Context) {
	var sysMenu map[string]interface{}
	if err := c.ShouldBindJSON(&sysMenu); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), sysMenu); err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询SysMenu
// @Tags     SysMenu
// @Summary  查询SysMenu
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询SysMenu"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysMenu/{id} [get]
func (cl *SysMenuHandle) Find(c *gin.Context) {
	var sysMenu domain.SysMenu
	if err := c.ShouldBindUri(&sysMenu); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), sysMenu)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询SysMenu列表
// @Tags     SysMenu
// @Summary  查询SysMenu列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageSysMenuSearch true "查询SysMenu列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysMenu/list [get]
func (cl *SysMenuHandle) List(c *gin.Context) {
	var sysMenu domain.PageSysMenuSearch
	if err := c.ShouldBindQuery(&sysMenu); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), sysMenu)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// RoleList 查询当前角色菜单列表
// @Tags     SysMenu
// @Summary  查询当前角色菜单列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysMenu/role-list [get]
func (cl *SysMenuHandle) RoleList(c *gin.Context) {
	rid, _ := common.GetRoleIdFromCtx(c)
	res, err := cl.s.RoleList(c.Request.Context(), rid)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// SetRoleList 设置角色菜单列表
// @Tags     SysMenu
// @Summary  设置角色菜单列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysRole true "设置角色菜单列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysMenu/role-list [put]
func (cl *SysMenuHandle) SetRoleList(c *gin.Context) {
	var sysRole domain.SysRole
	if err := c.ShouldBindJSON(&sysRole); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}
	err := cl.s.SetRoleList(c.Request.Context(), sysRole)
	if err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}
