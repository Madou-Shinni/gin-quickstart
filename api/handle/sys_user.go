package handle

import (
	"errors"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SysUserHandle struct {
	s *service.SysUserService
}

func NewSysUserHandle() *SysUserHandle {
	return &SysUserHandle{s: service.NewSysUserService()}
}

// Add 创建SysUser
// @Tags     SysUser
// @Summary  创建SysUser
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysUser true "创建SysUser"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysUser [post]
func (cl *SysUserHandle) Add(c *gin.Context) {
	var sysUser domain.SysUser
	if err := c.ShouldBindJSON(&sysUser); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(sysUser); err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, err.Error())
		return
	}

	response.Success(c)
}

// Delete 删除SysUser
// @Tags     SysUser
// @Summary  删除SysUser
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysUser true "删除SysUser"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysUser [delete]
func (cl *SysUserHandle) Delete(c *gin.Context) {
	var sysUser domain.SysUser
	if err := c.ShouldBindJSON(&sysUser); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(sysUser); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除SysUser
// @Tags     SysUser
// @Summary  批量删除SysUser
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除SysUser"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysUser/delete-batch [delete]
func (cl *SysUserHandle) DeleteByIds(c *gin.Context) {
	var ids request.Ids
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.DeleteByIds(ids); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Update 修改SysUser
// @Tags     SysUser
// @Summary  修改SysUser
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysUser true "修改SysUser"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysUser [put]
func (cl *SysUserHandle) Update(c *gin.Context) {
	var sysUser map[string]interface{}
	if err := c.ShouldBindJSON(&sysUser); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(sysUser); err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询SysUser
// @Tags     SysUser
// @Summary  查询SysUser
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询SysUser"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysUser/{id} [get]
func (cl *SysUserHandle) Find(c *gin.Context) {
	var sysUser domain.SysUser
	if err := c.ShouldBindUri(&sysUser); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(sysUser)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询SysUser列表
// @Tags     SysUser
// @Summary  查询SysUser列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageSysUserSearch true "查询SysUser列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysUser/list [get]
func (cl *SysUserHandle) List(c *gin.Context) {
	var sysUser domain.PageSysUserSearch
	if err := c.ShouldBindQuery(&sysUser); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(sysUser)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
