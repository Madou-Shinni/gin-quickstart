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

type SysCasbinHandle struct {
	s *service.SysCasbinService
}

func NewSysCasbinHandle() *SysCasbinHandle {
	return &SysCasbinHandle{s: service.NewSysCasbinService()}
}

// Add 创建SysCasbin
// @Tags     SysCasbin
// @Summary  设置角色权限
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.RolePermissionsReq true "创建SysCasbin"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysCasbin [post]
func (cl *SysCasbinHandle) Add(c *gin.Context) {
	var sysCasbin domain.RolePermissionsReq
	if err := c.ShouldBindJSON(&sysCasbin); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	err := cl.s.AddRolePermissions(c.Request.Context(), sysCasbin)
	if err != nil {
		response.Error(c, constant.CODE_ADD_FAILED)
		return
	}

	response.Success(c)
}

// Delete 删除SysCasbin
// @Tags     SysCasbin
// @Summary  删除SysCasbin
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysCasbin true "删除SysCasbin"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysCasbin [delete]
func (cl *SysCasbinHandle) Delete(c *gin.Context) {
	var sysCasbin domain.SysCasbin
	if err := c.ShouldBindJSON(&sysCasbin); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), sysCasbin); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除SysCasbin
// @Tags     SysCasbin
// @Summary  批量删除SysCasbin
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除SysCasbin"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysCasbin/delete-batch [delete]
func (cl *SysCasbinHandle) DeleteByIds(c *gin.Context) {
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

// Update 修改SysCasbin
// @Tags     SysCasbin
// @Summary  修改SysCasbin
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysCasbin true "修改SysCasbin"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysCasbin [put]
func (cl *SysCasbinHandle) Update(c *gin.Context) {
	var sysCasbin map[string]interface{}
	if err := c.ShouldBindJSON(&sysCasbin); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), sysCasbin); err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询SysCasbin
// @Tags     SysCasbin
// @Summary  查询SysCasbin
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询SysCasbin"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysCasbin/{id} [get]
func (cl *SysCasbinHandle) Find(c *gin.Context) {
	var sysCasbin domain.SysCasbin
	if err := c.ShouldBindUri(&sysCasbin); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), sysCasbin)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询SysCasbin列表
// @Tags     SysCasbin
// @Summary  查询SysCasbin列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageSysCasbinSearch true "查询SysCasbin列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysCasbin/list [get]
func (cl *SysCasbinHandle) List(c *gin.Context) {
	var sysCasbin domain.PageSysCasbinSearch
	if err := c.ShouldBindQuery(&sysCasbin); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), sysCasbin)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
