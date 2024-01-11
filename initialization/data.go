package initialization

import (
	"flag"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据初始化
// mysql、redis
func init() {
	ConfigInit()
	MysqlInit(conf.Conf.MysqlConfig)
	RedisInit(conf.Conf.RedisConfig)
}

// 初始化配置
// 将配置文件的信息反序列化到结构体中
func ConfigInit() {
	configFile := "configs/config.yml"
	s := flag.String("f", configFile, "choose config file.")
	flag.Parse()
	//viper.AddConfigPath(configPath)
	//viper.SetConfigName("config")     // 读取配置文件
	viper.SetConfigFile(*s)     // 读取配置文件
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() faild error:%v\n", err)
		return
	}
	// 把读取到的信息反序列化到Conf变量中
	if err := viper.Unmarshal(conf.Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
	}
	viper.WatchConfig()                            // （热加载时读取配置）监控配置文件
	viper.OnConfigChange(func(in fsnotify.Event) { // 配置文件修改时触发回调
		if err := viper.Unmarshal(conf.Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
		}
	})
}

// mysql连接初始化
func MysqlInit(config *conf.MysqlConfig) {
	// dsn := "root:123456@tcp(192.168.0.6:3306)/gin?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true, //打印sql
		//SkipDefaultTransaction: true, //禁用事务
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	// 自动迁移
	db.AutoMigrate(
		// 表
		domain.File{},
	)

	global.DB = db
}

// redis连接初始化
func RedisInit(config *conf.RedisConfig) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: 8,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	global.Rdb = rdb
	return
}

// 释放资源
func Close() {
	global.Rdb.Close()
}
