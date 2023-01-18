package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"gorm.io/gorm"
)

type Demo struct {
	gorm.Model
	Did  int64  `json:"did,omitempty" form:"did" gorm:"column:did;comment:唯一标识"`                                        // 唯一标识
	Uid  int64  `json:"uid,omitempty" form:"uid" gorm:"column:uid;index:idx_uid_code,unique;comment:用户id"`              // 用户id
	Code string `json:"code,omitempty" form:"code" gorm:"column:code;index:idx_uid_code,unique;comment:code值(验证请求来源);"` // code
	QPS  int64  `json:"qps,omitempty" form:"qps" gorm:"column:qps;comment:每秒请求量;"`                                      // 每秒请求量
}

type PageDemoSearch struct {
	Demo
	request.PageSearch
}

func (Demo) TableName() string {
	return "demo"
}
