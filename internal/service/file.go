package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/snowflake"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/upload"
	"github.com/Madou-Shinni/go-logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

// 定义接口
type FileRepo interface {
	Create(ctx context.Context, file domain.File) error
	Delete(ctx context.Context, file domain.File) error
	Update(ctx context.Context, file map[string]interface{}) error
	Find(ctx context.Context, file domain.File) (domain.File, error)
	List(ctx context.Context, page domain.PageFileSearch) ([]domain.File, int64, error)
	DeleteByIds(ctx context.Context, ids request.Ids) error
}

type FileService struct {
	repo FileRepo
}

func NewFileService() *FileService {
	return &FileService{repo: &data.FileRepo{}}
}

func (s *FileService) Add(ctx context.Context, file domain.File) error {
	// 1.生成唯一标识
	// 因为我们在全局初始化的时候已经初始化了雪花算法的机器节点
	// 所以我们可以直接使用
	id := snowflake.GenerateID()

	file.ID = id

	// 3.持久化入库
	if err := s.repo.Create(ctx, file); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(file)", zap.Error(err), zap.Any("domain.File", file))
		return err
	}

	return nil
}

func (s *FileService) Delete(ctx context.Context, file domain.File) error {
	if err := s.repo.Delete(ctx, file); err != nil {
		logger.Error("s.repo.Delete(file)", zap.Error(err), zap.Any("domain.File", file))
		return err
	}

	return nil
}

func (s *FileService) Update(ctx context.Context, file map[string]interface{}) error {
	if err := s.repo.Update(ctx, file); err != nil {
		logger.Error("s.repo.Update(file)", zap.Error(err), zap.Any("domain.File", file))
		return err
	}

	return nil
}

func (s *FileService) Find(ctx context.Context, file domain.File) (domain.File, error) {
	res, err := s.repo.Find(ctx, file)

	if err != nil {
		logger.Error("s.repo.Find(file)", zap.Error(err), zap.Any("domain.File", file))
		return res, err
	}

	return res, nil
}

func (s *FileService) List(ctx context.Context, page domain.PageFileSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, count, err := s.repo.List(ctx, page)
	if err != nil {
		logger.Error("s.repo.List(page)", zap.Error(err), zap.Any("domain.PageFileSearch", page))
		return pageRes, err
	}

	pageRes.List = data
	pageRes.Total = count

	return pageRes, nil
}

func (s *FileService) DeleteByIds(ctx context.Context, ids request.Ids) error {
	if err := s.repo.DeleteByIds(ctx, ids); err != nil {
		logger.Error("s.DeleteByIds(ids)", zap.Error(err), zap.Any("ids request.Ids", ids))
		return err
	}

	return nil
}

// 分片上传
func (s *FileService) UploadChunk(ctx context.Context, file domain.File, fileHeader *multipart.FileHeader) (domain.File, error) {
	var (
		err error
	)

	// 获取文件后缀
	fileSuffix := strings.Split(fileHeader.Filename, ".")[1]
	filePrefix := strconv.Itoa(file.Index)
	fileDir := strconv.FormatInt(file.ID, 10)

	// 查询分片是否已经上传
	flag := global.Rdb.HGet(
		fmt.Sprintf(constant.FileChunkHkey+strconv.FormatInt(file.ID, 10)),
		fmt.Sprintf(constant.FileChunkHFiled+strconv.Itoa(file.Index)),
	).Val()
	if flag == "1" {
		// 已经上传你直接返回
		file.IsFinish = true
		return file, nil
	}

	// 使用工具上传（目录 = ./static/fileschunk/fileDir/filePrefix.fileSuffix）
	dst := constant.FileChunkPathPrefix + fileDir + filePrefix + "." + fileSuffix
	if err = upload.Upload(fileHeader, dst); err != nil {
		logger.Error("upload.Upload(file)", zap.Error(err), zap.Any("file domain.File", file))
		return file, err
	} else {
		// redis 记录分片文件信息
		global.Rdb.HSet(
			fmt.Sprintf(constant.FileChunkHkey+strconv.FormatInt(file.ID, 10)),
			fmt.Sprintf(constant.FileChunkHFiled+strconv.Itoa(file.Index)),
			"1")

		file.IsFinish = true
		return file, nil
	}

}

// 秒传
func (s *FileService) UploadFast(ctx context.Context, file domain.File) (domain.File, error) {
	var (
		err error
	)

	// 秒传
	f, err := s.repo.Find(ctx, file)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 无法秒传
		return domain.File{}, err
	} else {
		// 成功秒传
		return f, nil
	}
}

// 合并分片
func (s *FileService) MergeChunk(ctx context.Context, file domain.File) (domain.File, error) {
	// 文件目录
	fileDir := strconv.FormatInt(file.ID, 10) + "/"
	dir := constant.FilePathPrefix + fileDir

	// 查询所有分片是否上传完成
	chunkMap := global.Rdb.HGetAll(
		fmt.Sprintf(constant.FileChunkHkey + strconv.FormatInt(file.ID, 10)),
	).Val()
	if len(chunkMap) < file.TotalChunk {
		// 分片未上传完成，无法合并
		logger.Warn("所有分片未上传完成，无法合并", zap.Int64("file.ID", file.ID))
		return domain.File{}, nil
	}
	// 合并文件
	dst, err := upload.MergeChunk(dir, file.FileName)
	if err != nil {
		logger.Error("upload.MergeChunk(file.ID, dir)", zap.Error(err))
		return domain.File{}, err
	}

	// 删除分片 ps: 可以用消息队列优化
	chunkDir := constant.FileChunkPathPrefix + fileDir
	err = os.RemoveAll(chunkDir)
	if err != nil {
		logger.Error("删除分片文件失败", zap.Error(err))
	}
	global.Rdb.Del(fmt.Sprintf(constant.FileChunkHkey + strconv.FormatInt(file.ID, 10)))

	// 返回文件地址
	file.FilePath = dst

	// 入库
	err = s.repo.Create(ctx, file)
	if err != nil {
		logger.Error("s.repo.Create(file)", zap.Error(err))
		return domain.File{}, err
	}

	return file, nil
}
