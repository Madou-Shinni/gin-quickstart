package test

import (
	"context"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
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

func TestSetNx(t *testing.T) {
	ctx := context.Background()

	err := db.WithContext(ctx).Create(&domain.Demo{Name: "ricardo", Age: 18, BirthDay: &model.LocalTime{Time: time.Now()}}).Error
	assert.Equal(t, nil, err)

	ok, err := rdb.SetNX(ctx, "demo", "ricardo", 0).Result()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ok)
	err = rdb.SetNX(ctx, "demo", "ricardo", time.Second*10).Err()
	assert.Equal(t, nil, err)
}

func TestCreateDemo(t *testing.T) {
	ctx := context.Background()

	var demoTestData = []domain.Demo{
		{
			Model: model.Model{ID: 1},
			Name:  "Alice",
			Age:   25,
			BirthDay: &model.LocalTime{
				Time: time.Date(1998, 4, 15, 0, 0, 0, 0, time.UTC),
			},
			Tags: datatypes.JSONSlice[string]{"golang", "backend", "gin"},
		},
		{
			Model: model.Model{ID: 2},
			Name:  "Bob",
			Age:   30,
			BirthDay: &model.LocalTime{
				Time: time.Date(1993, 7, 8, 0, 0, 0, 0, time.UTC),
			},
			Tags: datatypes.JSONSlice[string]{"frontend", "react", "typescript"},
		},
		{
			Model: model.Model{ID: 3},
			Name:  "Carol",
			Age:   22,
			BirthDay: &model.LocalTime{
				Time: time.Date(2001, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			Tags: datatypes.JSONSlice[string]{"design", "ui", "ux"},
		},
		{
			Model: model.Model{ID: 4},
			Name:  "David",
			Age:   28,
			BirthDay: &model.LocalTime{
				Time: time.Date(1996, 11, 5, 0, 0, 0, 0, time.UTC),
			},
			Tags: datatypes.JSONSlice[string]{"golang", "docker", "kubernetes"},
		},
		{
			Model: model.Model{ID: 5},
			Name:  "Eve",
			Age:   35,
			BirthDay: &model.LocalTime{
				Time: time.Date(1989, 6, 30, 0, 0, 0, 0, time.UTC),
			},
			Tags: datatypes.JSONSlice[string]{"devops", "ci/cd"},
		},
	}

	err := db.WithContext(ctx).Create(&demoTestData).Error
	assert.Equal(t, nil, err)
}
