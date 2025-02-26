package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type Demo struct {
	model.Model
	Name     string           `json:"name"`
	Age      int              `json:"age"`
	BirthDay *model.LocalTime `json:"birth_day"`
}

type PageDemoSearch struct {
	Demo
	request.PageSearch
}

func (Demo) TableName() string {
	return "demo"
}
