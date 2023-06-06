package handle

import (
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/gin-gonic/gin"
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
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(demo); err != nil {
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

	if err := cl.s.Delete(demo); err != nil {
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

	if err := cl.s.DeleteByIds(ids); err != nil {
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
	var demo map[string]interface{}
	if err := c.ShouldBindJSON(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(demo); err != nil {
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
// @Param    data query     domain.Demo true "查询Demo"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /demo [get]
func (cl *DemoHandle) Find(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindQuery(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(demo)

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

	res, err := cl.s.List(demo)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
