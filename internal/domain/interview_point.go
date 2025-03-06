package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
)

type InterviewPoint struct {
	model.Model
	Name       string  `json:"name"`
	Latitude   float64 `json:"latitude"`    // 纬度
	Longitude  float64 `json:"longitude"`   // 经度
	RiderCount int     `json:"rider_count"` // 当前骑手数量
}

type PageInterviewPointSearch struct {
	InterviewPoint
	request.PageSearch
}

func (InterviewPoint) TableName() string {
	return "interview_point"
}
