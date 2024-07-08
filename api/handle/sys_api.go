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

type SysApiHandle struct {
	s *service.SysApiService
}

func NewSysApiHandle() *SysApiHandle {
	return &SysApiHandle{s: service.NewSysApiService()}
}

// Add 创建SysApi
// @Tags     SysApi
// @Summary  创建SysApi
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysApi true "创建SysApi"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysApi [post]
func (cl *SysApiHandle) Add(c *gin.Context) {
	var sysApi domain.SysApi
	if err := c.ShouldBindJSON(&sysApi); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), sysApi); err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除SysApi
// @Tags     SysApi
// @Summary  删除SysApi
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysApi true "删除SysApi"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysApi [delete]
func (cl *SysApiHandle) Delete(c *gin.Context) {
	var sysApi domain.SysApi
	if err := c.ShouldBindJSON(&sysApi); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), sysApi); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除SysApi
// @Tags     SysApi
// @Summary  批量删除SysApi
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除SysApi"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysApi/delete-batch [delete]
func (cl *SysApiHandle) DeleteByIds(c *gin.Context) {
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

// Update 修改SysApi
// @Tags     SysApi
// @Summary  修改SysApi
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.SysApi true "修改SysApi"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /sysApi [put]
func (cl *SysApiHandle) Update(c *gin.Context) {
	var sysApi map[string]interface{}
	if err := c.ShouldBindJSON(&sysApi); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), sysApi); err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询SysApi
// @Tags     SysApi
// @Summary  查询SysApi
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询SysApi"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysApi/{id} [get]
func (cl *SysApiHandle) Find(c *gin.Context) {
	var sysApi domain.SysApi
	if err := c.ShouldBindUri(&sysApi); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), sysApi)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询SysApi列表
// @Tags     SysApi
// @Summary  查询SysApi列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageSysApiSearch true "查询SysApi列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /sysApi/list [get]
func (cl *SysApiHandle) List(c *gin.Context) {
	var sysApi domain.PageSysApiSearch
	if err := c.ShouldBindQuery(&sysApi); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), sysApi)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
