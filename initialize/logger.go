package initialize

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/go-logger"
	"time"
)

// 自定义日志初始化配置
func init() {
	var err error

	env := conf.Conf.Env
	file := fmt.Sprint(conf.Conf.LogFile)

	if env == "prod" {
		// 生产环境
		err = initProd(file)
		if err != nil {
			panic(err)
		}
		return
	}

	err = initdev()
	if err != nil {
		panic(err)
	}
}

// 开发日志初始化配置
func initdev() error {
	var err error

	_, err = logger.NewJSONLogger(
		// 日志等级
		logger.WithDebugLevel(),
		// 时间格式化
		logger.WithTimeLayout(time.DateTime),
	)

	return err
}

// 生产日志初始化配置
func initProd(file string) error {
	_, err := logger.NewJSONLogger(
		// 日志等级
		logger.WithDebugLevel(),
		// 写出的文件
		logger.WithFileRotationP(file),
		// 不在控制台打印
		logger.WithDisableConsole(),
		// 时间格式化
		logger.WithTimeLayout(time.DateTime),
	)

	return err
}
