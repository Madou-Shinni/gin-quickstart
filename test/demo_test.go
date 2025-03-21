package test

import (
	"context"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestDemo(t *testing.T) {
	ctx := context.Background()

	err := db.WithContext(ctx).Create(&domain.Demo{Name: "ricardo", Age: 18, BirthDay: &model.LocalTime{Time: time.Now()}}).Error
	assert.Equal(t, nil, err)

	err = rdb.Set(ctx, "demo", "ricardo", time.Second*10).Err()
	assert.Equal(t, nil, err)
}
