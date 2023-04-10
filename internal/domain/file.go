package domain

import "github.com/Madou-Shinni/gin-quickstart/pkg/request"

type File struct {
	ID           int64  `json:"id,omitempty" form:"id" gorm:"column:id;comment:主键;primarykey"`                          // 文件唯一标识
	FileMd5      string `json:"fileMd5,omitempty" form:"fileMd5" gorm:"column:file_md5;comment:文件MD5"`                  // 文件MD5
	FileSize     string `json:"fileSize,omitempty" form:"fileSize" gorm:"column:file_size;comment:文件大小"`                // 文件大小
	FilePath     string `json:"filePath,omitempty" form:"filePath" gorm:"column:file_path;comment:文件路径"`                // 文件路径
	FileName     string `json:"fileName,omitempty" form:"fileName" gorm:"column:file_name;comment:文件名"`                 // 文件名
	TotalChunk   int    `json:"totalChunk,omitempty" form:"totalChunk" gorm:"column:total_chunk;comment:文件总分片数"`        // 文件总分片数
	AlreadyChunk string `json:"alreadyChunk,omitempty" form:"alreadyChunk" gorm:"column:already_chunk;comment:已经上传的分片"` // 已经上传的分片
	Index        int    `json:"index,omitempty" form:"index" gorm:"-"`                                                  // 当前分片
	IsFinish     bool   `json:"isFinish,omitempty" form:"isFinish" gorm:"-"`                                            // 是否完成
}

type PageFileSearch struct {
	File
	request.PageSearch
}

func (File) TableName() string {
	return "file"
}
