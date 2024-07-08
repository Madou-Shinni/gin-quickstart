package handle

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/snowflake"
	"github.com/Madou-Shinni/go-logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type FileHandle struct {
	s *service.FileService
}

func NewFileHandle() *FileHandle {
	return &FileHandle{s: service.NewFileService()}
}

// Add 创建File
// @Tags     File
// @Summary  创建File
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.File true "创建File"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /file [post]
func (cl *FileHandle) Add(c *gin.Context) {
	var file domain.File
	if err := c.ShouldBindJSON(&file); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(c.Request.Context(), file); err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// MergeChunk 合并分片文件
// @Tags     File
// @Summary  合并分片文件
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.File true "合并分片文件"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /file/merge-chunk [get]
func (cl *FileHandle) MergeChunk(c *gin.Context) {
	var (
		file domain.File
		err  error
	)

	if err := c.ShouldBindQuery(&file); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if file, err = cl.s.MergeChunk(c.Request.Context(), file); err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c, file)
}

// Upload 普通上传File
// @Tags     File
// @Summary  普通上传File
// @accept   multipart/form-data
// @Produce  application/json
// @Param    file formData     file true "普通上传File"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /file/upload [post]
func (cl *FileHandle) Upload(c *gin.Context) {
	var (
		err      error
		filePath string
	)

	form, _ := c.MultipartForm()
	fileHeaders := form.File["file"]

	for _, fileHeader := range fileHeaders {
		suffix := filepath.Ext(fileHeader.Filename)
		filePath = fmt.Sprint(conf.Conf.UploadConfig.Dir, "/", uuid.NewString(), suffix)
		os.MkdirAll(conf.Conf.UploadConfig.Dir, os.ModePerm)
		err = c.SaveUploadedFile(fileHeader, filePath)
		if err != nil {
			logger.Error("SaveUploadedFile", zap.Error(err))
			response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
			return
		}
	}

	response.Success(c, filePath)
}

// UploadChunk 分片上传File
// @Tags     File
// @Summary  分片上传File
// @accept   multipart/form-data
// @Produce  application/json
// @Param    file formData     file true "分片上传File"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /file/upload-chunk [post]
func (cl *FileHandle) UploadChunk(c *gin.Context) {
	var (
		file domain.File
		err  error
	)

	form, _ := c.MultipartForm()
	fileHeaders := form.File["file"]

	if err := c.ShouldBindQuery(&file); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	for _, fileHeader := range fileHeaders {
		file, err = cl.s.UploadChunk(c.Request.Context(), file, fileHeader)
		if err != nil {
			response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
			return
		}
	}

	response.Success(c, file)
}

// Chunkid 获取分块文件id
// @Tags     File
// @Summary  获取分块文件id
// @accept   application/json
// @Produce  application/json
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /file/chunkid [get]
func (cl *FileHandle) Chunkid(c *gin.Context) {
	var (
		file domain.File
	)

	file.ID = snowflake.GenerateID()

	response.Success(c, file)
}

// Delete 删除File
// @Tags     File
// @Summary  删除File
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.File true "删除File"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /file [delete]
func (cl *FileHandle) Delete(c *gin.Context) {
	var file domain.File
	if err := c.ShouldBindJSON(&file); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(c.Request.Context(), file); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// DeleteByIds 批量删除File
// @Tags     File
// @Summary  批量删除File
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.Ids true "批量删除File"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /file/delete-batch [delete]
func (cl *FileHandle) DeleteByIds(c *gin.Context) {
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

// Update 修改File
// @Tags     File
// @Summary  修改File
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.File true "修改File"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /file [put]
func (cl *FileHandle) Update(c *gin.Context) {
	var file map[string]interface{}
	if err := c.ShouldBindJSON(&file); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(c.Request.Context(), file); err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询File
// @Tags     File
// @Summary  查询File
// @accept   application/json
// @Produce  application/json
// @Param    fileMd5 query     string true "查询File"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /file [get]
func (cl *FileHandle) Find(c *gin.Context) {
	var file domain.File
	if err := c.ShouldBindQuery(&file); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(c.Request.Context(), file)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询File列表
// @Tags     File
// @Summary  查询File列表
// @accept   application/json
// @Produce  application/json
// @Param    data query     domain.File true "查询File列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /file/list [get]
func (cl *FileHandle) List(c *gin.Context) {
	var file domain.PageFileSearch
	if err := c.ShouldBindQuery(&file); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(c.Request.Context(), file)

	if err != nil {
		response.Error(c, constant.CODE_FIND_FAILED, constant.CODE_FIND_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
