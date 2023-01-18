package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// 全局变量
var (
	DB  *gorm.DB
	Rdb *redis.Client
)
