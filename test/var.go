package test

import (
	"context"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	dbConfig = struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "123456",
		Database: "test",
	}

	rdbConfig = struct {
		Addr     string
		Password string
		DB       int
	}{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}

	db  *global.Data
	rdb *redis.Client
)

func setup() {
	db = getDB()
	rdb = getRdb()
	log.Println("> setup completed")
}

func teardown() {
	rdb.Close()
	log.Println("> teardown completed")
}

func getDB() *global.Data {
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields:                              true, //打印sql
		DisableForeignKeyConstraintWhenMigrating: true, //禁用外键约束
		//SkipDefaultTransaction: true, //禁用事务
	})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := d.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)

	return global.NewData(d)
}

func getRdb() *redis.Client {
	ctx := context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     rdbConfig.Addr,
		Password: rdbConfig.Password,
		DB:       rdbConfig.DB,
		PoolSize: 8,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return rdb
}
