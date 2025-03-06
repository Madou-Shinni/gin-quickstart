package test

import (
	"context"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 约束数值类型
type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// 泛型函数：接收数值类型，返回指针
func ToPointer[T Number](value T) *T {
	return &value
}

func TestPush(t *testing.T) {
	ctx := context.Background()
	rlService := service.NewResumeLogService()
	var list = []service.ResumePushReq{
		{
			RID:         1,
			SelectedIID: ToPointer(uint(1)),
			Latitude:    0,
			Longitude:   0,
		},
		{
			RID:         2,
			SelectedIID: nil,
			Latitude:    31.2305,
			Longitude:   121.4800,
		},
		{
			RID:         3,
			SelectedIID: nil,
			Latitude:    0,
			Longitude:   0,
		},
		{
			RID:         4,
			SelectedIID: nil,
			Latitude:    31.1700,
			Longitude:   121.4450,
		},
		{
			RID:         5,
			SelectedIID: nil,
			Latitude:    0,
			Longitude:   0,
		},
	}
	for _, req := range list {
		ip, err := rlService.Push(ctx, req)
		assert.Equal(t, nil, err)
		t.Logf("ip:%v", ip)
	}
}
