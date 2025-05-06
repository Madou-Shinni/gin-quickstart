package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"gorm.io/datatypes"
)

type Demo struct {
	model.Model
	Name     string                      `json:"name" form:"name"`
	Age      int                         `json:"age"`
	BirthDay *model.LocalTime            `json:"birth_day"`
	Tags     datatypes.JSONSlice[string] `json:"tags"`
}

type PageDemoSearch struct {
	Demo
	request.PageSearch
	TagsQuery string `json:"tagsQuery" form:"tagsQuery"`
}

func (Demo) TableName() string {
	return "demo"
}
