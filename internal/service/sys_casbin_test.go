package service

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var (
	user        = "root"
	pwd         = "123456"
	host        = "localhost"
	port        = "3306"
	dbName      = "casbin"
	maxIdleConn = 10
	maxOpenConn = 10
)

func TestCasbin(t *testing.T) {
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		user, pwd, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true, //打印sql
		//SkipDefaultTransaction: true, //禁用事务
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(maxIdleConn)
	sqlDB.SetMaxOpenConns(maxOpenConn)
	if err != nil {
		panic(err)
	}

	global.DB = db
	service := NewSysCasbinService()
	t.Log(service.e)
	s2 := NewSysCasbinService()
	t.Log(s2.e)
}
