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

type SysRoleHandle struct {
	s *service.SysRoleService
}

func NewSysRoleHandle() *SysRoleHandle {
	return &SysRoleHandle{s: service.NewSysRoleService()}
}

// Add 创建SysRole
// @Tags     SysRole
// @Summary  创建SysRole
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysRole true "创建SysRole"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysRole [post]
func (cl *SysRoleHandle) Add(c *gin.Context) {
	var sysRole domain.SysRole
	if err := c.ShouldBindJSON(&sysRole); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), sysRole); err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除SysRole
// @Tags     SysRole
// @Summary  删除SysRole
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysRole true "删除SysRole"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysRole [delete]
func (cl *SysRoleHandle) Delete(c *gin.Context) {
	var sysRole domain.SysRole
	if err := c.ShouldBindJSON(&sysRole); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), sysRole); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除SysRole
// @Tags     SysRole
// @Summary  批量删除SysRole
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除SysRole"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysRole/delete-batch [delete]
func (cl *SysRoleHandle) DeleteByIds(c *gin.Context) {
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

// Update 修改SysRole
// @Tags     SysRole
// @Summary  修改SysRole
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysRole true "修改SysRole"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysRole [put]
func (cl *SysRoleHandle) Update(c *gin.Context) {
	var sysRole map[string]interface{}
	if err := c.ShouldBindJSON(&sysRole); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), sysRole); err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询SysRole
// @Tags     SysRole
// @Summary  查询SysRole
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询SysRole"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysRole/{id} [get]
func (cl *SysRoleHandle) Find(c *gin.Context) {
	var sysRole domain.SysRole
	if err := c.ShouldBindUri(&sysRole); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), sysRole)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询SysRole列表
// @Tags     SysRole
// @Summary  查询SysRole列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageSysRoleSearch true "查询SysRole列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysRole/list [get]
func (cl *SysRoleHandle) List(c *gin.Context) {
	var sysRole domain.PageSysRoleSearch
	if err := c.ShouldBindQuery(&sysRole); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), sysRole)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// SetUserRoleList 设置用户角色列表
// @Tags     SysRole
// @Summary  设置用户角色列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysUser true "设置用户角色列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysRole/user-list [put]
func (cl *SysRoleHandle) SetUserRoleList(c *gin.Context) {
	var SysUser domain.SysUser
	if err := c.ShouldBindJSON(&SysUser); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	SysUser.DefaultRole, _ = common.GetRoleIdFromCtx(c)

	err := cl.s.SetUserRoleList(c.Request.Context(), SysUser)
	if err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}
