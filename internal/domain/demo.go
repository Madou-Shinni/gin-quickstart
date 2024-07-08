package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type Demo struct {
	model.Model
	//F string `json:"f" binding:"required"`
	//N string `json:"n" binding:"required"`
}

type PageDemoSearch struct {
	Demo
	request.PageSearch
}

func (Demo) TableName() string {
	return "demo"
}
