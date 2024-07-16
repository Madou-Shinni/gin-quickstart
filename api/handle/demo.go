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

type DemoHandle struct {
	s *service.DemoService
}

func NewDemoHandle() *DemoHandle {
	return &DemoHandle{s: service.NewDemoService()}
}

// Add 创建Demo
// @Tags     Demo
// @Summary  创建Demo
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.Demo true "创建Demo"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /demo [post]
func (cl *DemoHandle) Add(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindJSON(&demo); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), demo); err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除Demo
// @Tags     Demo
// @Summary  删除Demo
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.Demo true "删除Demo"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /demo [delete]
func (cl *DemoHandle) Delete(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindJSON(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), demo); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除Demo
// @Tags     Demo
// @Summary  批量删除Demo
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.Ids true "批量删除Demo"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /demo/delete-batch [delete]
func (cl *DemoHandle) DeleteByIds(c *gin.Context) {
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

// Update 修改Demo
// @Tags     Demo
// @Summary  修改Demo
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.Demo true "修改Demo"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /demo [put]
func (cl *DemoHandle) Update(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindJSON(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), demo); err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询Demo
// @Tags     Demo
// @Summary  查询Demo
// @accept   application/json
// @Produce  application/json
// @Param    id path     uint true "查询Demo"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /demo/{id} [get]
func (cl *DemoHandle) Find(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindUri(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), demo)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询Demo列表
// @Tags     Demo
// @Summary  查询Demo列表
// @accept   application/json
// @Produce  application/json
// @Param    data query     domain.Demo true "查询Demo列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /demo/list [get]
func (cl *DemoHandle) List(c *gin.Context) {
	var demo domain.PageDemoSearch
	if err := c.ShouldBindQuery(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), demo)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
