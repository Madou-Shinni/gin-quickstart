package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type DataImport struct {
	model.Model
}

type PageDataImportSearch struct {
	DataImport
	request.PageSearch
}

func (DataImport) TableName() string {
	return "data_import"
}
