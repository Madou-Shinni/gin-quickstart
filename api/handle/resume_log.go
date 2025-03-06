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

type ResumeLogHandle struct {
	s *service.ResumeLogService
}

func NewResumeLogHandle() *ResumeLogHandle {
	return &ResumeLogHandle{s: service.NewResumeLogService()}
}

// Add 创建ResumeLog
// @Tags     ResumeLog
// @Summary  创建ResumeLog
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.ResumeLog true "创建ResumeLog"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /resumeLog [post]
func (cl *ResumeLogHandle) Add(c *gin.Context) {
	var resumeLog domain.ResumeLog
	if err := c.ShouldBindJSON(&resumeLog); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), resumeLog); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除ResumeLog
// @Tags     ResumeLog
// @Summary  删除ResumeLog
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.ResumeLog true "删除ResumeLog"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /resumeLog [delete]
func (cl *ResumeLogHandle) Delete(c *gin.Context) {
	var resumeLog domain.ResumeLog
	if err := c.ShouldBindJSON(&resumeLog); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), resumeLog); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除ResumeLog
// @Tags     ResumeLog
// @Summary  批量删除ResumeLog
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除ResumeLog"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /resumeLog/delete-batch [delete]
func (cl *ResumeLogHandle) DeleteByIds(c *gin.Context) {
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

// Update 修改ResumeLog
// @Tags     ResumeLog
// @Summary  修改ResumeLog
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.ResumeLog true "修改ResumeLog"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /resumeLog [put]
func (cl *ResumeLogHandle) Update(c *gin.Context) {
	var resumeLog domain.ResumeLog
	if err := c.ShouldBindJSON(&resumeLog); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), resumeLog); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询ResumeLog
// @Tags     ResumeLog
// @Summary  查询ResumeLog
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询ResumeLog"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /resumeLog/{id} [get]
func (cl *ResumeLogHandle) Find(c *gin.Context) {
	var resumeLog domain.ResumeLog
	if err := c.ShouldBindUri(&resumeLog); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), resumeLog)

	if err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询ResumeLog列表
// @Tags     ResumeLog
// @Summary  查询ResumeLog列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageResumeLogSearch true "查询ResumeLog列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /resumeLog/list [get]
func (cl *ResumeLogHandle) List(c *gin.Context) {
	var resumeLog domain.PageResumeLogSearch
	if err := c.ShouldBindQuery(&resumeLog); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), resumeLog)

	if err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

func (cl *ResumeLogHandle) Push(c *gin.Context) {
	var resumeLog service.ResumePushReq
	if err := c.ShouldBindJSON(&resumeLog); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if res, err := cl.s.Push(c.Request.Context(), resumeLog); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	} else {
		response.Success(c, gin.H{
			"ip":  res,
			"rid": resumeLog.RID,
		})
	}
}
