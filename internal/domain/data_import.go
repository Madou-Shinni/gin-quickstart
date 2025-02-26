package domain

import (
	"encoding/json"
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"gorm.io/datatypes"
)

type FailedReason struct {
	Row    string `gorm:"type:varchar(128);default:'';" json:"row"`
	Reason string `gorm:"type:varchar(128);default:'';comment:'失败原因'" json:"reason"`
}

type DataImport struct {
	model.Model
	FileName string `gorm:"type:varchar(64);not null;" json:"filename"`
	FileUrl  string `gorm:"type:varchar(512);not null;" json:"file_url"`
	// importing： 导入中，success：导入成功，导入失败：failed
	Status        string                            `gorm:"type:varchar(16);index;default:importing;not null;" json:"status" form:"status"`
	Category      string                            `gorm:"type:varchar(32);index;default:'';not null;" json:"category" form:"category"`
	Count         uint                              `gorm:"type:int;default:0;not null;" json:"count"`
	SuccessCount  uint                              `gorm:"type:int;default:0;not null;" json:"success_count"`
	FailureCount  uint                              `gorm:"type:int;default:0;not null;" json:"failure_count"`
	FailedReasons datatypes.JSONSlice[FailedReason] `gorm:"type:json" json:"failed_reasons"`    // 错误信息
	Data          json.RawMessage                   `json:"data" gorm:"-" swaggerignore:"true"` // 导入数据
}

type PageDataImportSearch struct {
	DataImport
	request.PageSearch
}

func (DataImport) TableName() string {
	return "data_import"
}
