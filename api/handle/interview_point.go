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

type InterviewPointHandle struct {
	s *service.InterviewPointService
}

func NewInterviewPointHandle() *InterviewPointHandle {
	return &InterviewPointHandle{s: service.NewInterviewPointService()}
}

// Add 创建InterviewPoint
// @Tags     InterviewPoint
// @Summary  创建InterviewPoint
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.InterviewPoint true "创建InterviewPoint"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /interviewPoint [post]
func (cl *InterviewPointHandle) Add(c *gin.Context) {
	var interviewPoint domain.InterviewPoint
	if err := c.ShouldBindJSON(&interviewPoint); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), interviewPoint); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除InterviewPoint
// @Tags     InterviewPoint
// @Summary  删除InterviewPoint
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.InterviewPoint true "删除InterviewPoint"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /interviewPoint [delete]
func (cl *InterviewPointHandle) Delete(c *gin.Context) {
	var interviewPoint domain.InterviewPoint
	if err := c.ShouldBindJSON(&interviewPoint); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), interviewPoint); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除InterviewPoint
// @Tags     InterviewPoint
// @Summary  批量删除InterviewPoint
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除InterviewPoint"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /interviewPoint/delete-batch [delete]
func (cl *InterviewPointHandle) DeleteByIds(c *gin.Context) {
	var ids request.Ids
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.DeleteByIds(c.Request.Context(), ids); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Update 修改InterviewPoint
// @Tags     InterviewPoint
// @Summary  修改InterviewPoint
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.InterviewPoint true "修改InterviewPoint"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /interviewPoint [put]
func (cl *InterviewPointHandle) Update(c *gin.Context) {
	var interviewPoint domain.InterviewPoint
	if err := c.ShouldBindJSON(&interviewPoint); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), interviewPoint); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询InterviewPoint
// @Tags     InterviewPoint
// @Summary  查询InterviewPoint
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询InterviewPoint"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /interviewPoint/{id} [get]
func (cl *InterviewPointHandle) Find(c *gin.Context) {
	var interviewPoint domain.InterviewPoint
	if err := c.ShouldBindUri(&interviewPoint); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), interviewPoint)

	if err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询InterviewPoint列表
// @Tags     InterviewPoint
// @Summary  查询InterviewPoint列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageInterviewPointSearch true "查询InterviewPoint列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /interviewPoint/list [get]
func (cl *InterviewPointHandle) List(c *gin.Context) {
	var interviewPoint domain.PageInterviewPointSearch
	if err := c.ShouldBindQuery(&interviewPoint); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), interviewPoint)

	if err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
