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

type DataImportHandle struct {
	s *service.DataImportService
}

func NewDataImportHandle() *DataImportHandle {
	return &DataImportHandle{s: service.NewDataImportService()}
}

// Add 创建DataImport
// @Tags     DataImport
// @Summary  创建DataImport
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.DataImport true "创建DataImport"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /dataImport [post]
func (cl *DataImportHandle) Add(c *gin.Context) {
	var dataImport domain.DataImport
	if err := c.ShouldBindJSON(&dataImport); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			response.Error(c, constant.CODE_INVALID_PARAMETER, tools.TransErrs(errs))
			return
		}
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), dataImport); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除DataImport
// @Tags     DataImport
// @Summary  删除DataImport
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.DataImport true "删除DataImport"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /dataImport [delete]
func (cl *DataImportHandle) Delete(c *gin.Context) {
	var dataImport domain.DataImport
	if err := c.ShouldBindJSON(&dataImport); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), dataImport); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除DataImport
// @Tags     DataImport
// @Summary  批量删除DataImport
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     request.Ids true "批量删除DataImport"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /dataImport/delete-batch [delete]
func (cl *DataImportHandle) DeleteByIds(c *gin.Context) {
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

// Update 修改DataImport
// @Tags     DataImport
// @Summary  修改DataImport
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.DataImport true "修改DataImport"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /dataImport [put]
func (cl *DataImportHandle) Update(c *gin.Context) {
	var dataImport domain.DataImport
	if err := c.ShouldBindJSON(&dataImport); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), dataImport); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询DataImport
// @Tags     DataImport
// @Summary  查询DataImport
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    id path     uint true "查询DataImport"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /dataImport/{id} [get]
func (cl *DataImportHandle) Find(c *gin.Context) {
	var dataImport domain.DataImport
	if err := c.ShouldBindUri(&dataImport); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), dataImport)

	if err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询DataImport列表
// @Tags     DataImport
// @Summary  查询DataImport列表
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     domain.PageDataImportSearch true "查询DataImport列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /dataImport/list [get]
func (cl *DataImportHandle) List(c *gin.Context) {
	var dataImport domain.PageDataImportSearch
	if err := c.ShouldBindQuery(&dataImport); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), dataImport)

	if err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// Template 模板
// @Tags     DataImport
// @Summary  模板
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     service.DataImportTemplateReq true "类型"
// @Success  200  {string} string            "{"code":200,"msg":"模板","data":{}"}"
// @Router   /dataImport/template [get]
func (cl *DataImportHandle) Template(c *gin.Context) {
	var req service.DataImportTemplateReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Template(c.Request.Context(), req)

	if err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, err.Error())
		return
	}

	response.Success(c, res)
}

// Import 导入数据
// @Tags     DataImport
// @Summary  导入数据
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data body     domain.DataImport true "导入数据"
// @Success  200  {string} string            "{"code":200,"msg":"导入数据成功","data":{}"}"
// @Router   /dataImport/import [post]
func (cl *DataImportHandle) Import(c *gin.Context) {
	var DataImport domain.DataImport
	if err := c.ShouldBindJSON(&DataImport); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, err.Error())
		return
	}

	res, err := cl.s.Import(c.Request.Context(), DataImport)

	if err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, err.Error())
		return
	}

	response.Success(c, res)
}
