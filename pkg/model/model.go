package model

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type LocalTime struct {
	time.Time
}

type LocalDate struct {
	time.Time
}

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id" uri:"id"`         // 主键
	CreatedAt *LocalTime     `json:"createdAt" form:"createdAt" swaggerignore:"true"` // 创建时间
	UpdatedAt *LocalTime     `json:"updatedAt" form:"updatedAt" swaggerignore:"true"` // 修改时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}

	if !strings.Contains(string(data), "T") {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), loc)
		if err != nil {
			return err
		}
		*t = LocalTime{Time: now}
	} else {
		parse, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
		if err != nil {
			return err
		}
		*t = LocalTime{Time: parse}
	}

	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)

}

func (t LocalDate) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))
	return []byte(formatted), nil
}

func (t LocalDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *LocalDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)

}
