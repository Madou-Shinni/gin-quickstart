package initialize

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	defaultConfigPath = "./configs"
)

// 数据初始化
// mysql、redis
func init() {
	flag.String("c", defaultConfigPath, "choose config file.")
	ConfigInit()
	MysqlInit(conf.Conf.MysqlConfig)
	RedisInit(conf.Conf.RedisConfig)
}

// 初始化配置
// 将配置文件的信息反序列化到结构体中
func ConfigInit() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	configPath := viper.Get("c").(string)
	viper.AddConfigPath(configPath)

	// 获取目录下所有的文件
	files, err := getFilesInDir(configPath)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	// 根据文件名排序，以确保按照指定顺序读取配置文件
	sortedFiles := sortFiles(files)

	// 按顺序读取每个配置文件
	for _, file := range sortedFiles {
		fileName := filepath.Base(file)
		viper.SetConfigName(fileName[:len(fileName)-len(filepath.Ext(fileName))])

		// 合并读取配置文件，忽略不存在的文件
		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				log.Fatalln("Config file not found")
			} else {
				// Config file was found but another error was produced
				log.Fatalf("Config file err: %v", err)
			}
		}

		log.Printf("Successfully read %s\n", fileName)
	}

	// 把读取到的信息反序列化到Conf变量中
	if err := viper.Unmarshal(conf.Conf); err != nil {
		log.Printf("viper.Unmarshal failed,err:%v\n", err)
	}
	viper.WatchConfig()                            // （热加载时读取配置）监控配置文件
	viper.OnConfigChange(func(in fsnotify.Event) { // 配置文件修改时触发回调
		if err := viper.Unmarshal(conf.Conf); err != nil {
			log.Printf("viper.Unmarshal failed,err:%v\n", err)
		}
		log.Println("Config file changed")
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

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	if err != nil {
		log.Println(err)
		return
	}

	// 自动迁移
	db.AutoMigrate(
		// 表
		domain.File{},
		domain.SysUser{},
		domain.SysRole{},
		//domain.SysCasbin{},
		domain.SysApi{},
		domain.SysMenu{},
	)

	global.DB = db

	//plugin := gorm_plugin.NewLogPlugin()
	//plugin.Apply(global.DB)
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
		log.Println(err)
		return
	}

	global.Rdb = rdb
	return
}

// 释放资源
func Close() {
	if global.Rdb == nil {
		return
	}
	global.Rdb.Close()
}

// 获取目录下所有的yaml文件
func getFilesInDir(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(info.Name())[1:] == "yaml" || filepath.Ext(info.Name())[1:] == "yml") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// 根据文件名排序文件列表
func sortFiles(files []string) []string {
	// 以文件名排序
	// 这里可以根据自己的需求定义更复杂的排序规则
	sortedFiles := files
	for i := 0; i < len(sortedFiles)-1; i++ {
		for j := i + 1; j < len(sortedFiles); j++ {
			if filepath.Base(sortedFiles[i]) > filepath.Base(sortedFiles[j]) {
				sortedFiles[i], sortedFiles[j] = sortedFiles[j], sortedFiles[i]
			}
		}
	}
	return sortedFiles
}
