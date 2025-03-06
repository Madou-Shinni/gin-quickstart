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

type ResumeHandle struct {
	s *service.ResumeService
}

func NewResumeHandle() *ResumeHandle {
	return &ResumeHandle{s: service.NewResumeService()}
}

// Add 创建Resume
// @Tags     Resume
// @Summary  创建Resume
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.Resume true "创建Resume"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /resume [post]
func (cl *ResumeHandle) Add(c *gin.Context) {
	var resume domain.Resume
	if err := c.ShouldBindJSON(&resume); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), resume); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除Resume
// @Tags     Resume
// @Summary  删除Resume
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.Resume true "删除Resume"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /resume [delete]
func (cl *ResumeHandle) Delete(c *gin.Context) {
	var resume domain.Resume
	if err := c.ShouldBindJSON(&resume); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), resume); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除Resume
// @Tags     Resume
// @Summary  批量删除Resume
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除Resume"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /resume/delete-batch [delete]
func (cl *ResumeHandle) DeleteByIds(c *gin.Context) {
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

// Update 修改Resume
// @Tags     Resume
// @Summary  修改Resume
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.Resume true "修改Resume"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /resume [put]
func (cl *ResumeHandle) Update(c *gin.Context) {
	var resume domain.Resume
	if err := c.ShouldBindJSON(&resume); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), resume); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询Resume
// @Tags     Resume
// @Summary  查询Resume
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询Resume"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /resume/{id} [get]
func (cl *ResumeHandle) Find(c *gin.Context) {
	var resume domain.Resume
	if err := c.ShouldBindUri(&resume); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), resume)

	if err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询Resume列表
// @Tags     Resume
// @Summary  查询Resume列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageResumeSearch true "查询Resume列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /resume/list [get]
func (cl *ResumeHandle) List(c *gin.Context) {
	var resume domain.PageResumeSearch
	if err := c.ShouldBindQuery(&resume); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), resume)

	if err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
