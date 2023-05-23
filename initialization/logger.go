package initialization

import (
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/go-logger"
	"strconv"
	"time"
)

// 自定义日志初始化配置
func init() {
	var err error

	env := conf.Conf.Env
	file := "./log/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".log"

	if env == "dev" {
		// 开发环境
		_, err = logger.NewJSONLogger(
			// 日志等级
			logger.WithDebugLevel(),
			// 时间格式化
			logger.WithTimeLayout("2006-01-02 15:04:05"),
		)
	} else if env == "prod" {
		// 生产环境
		_, err = logger.NewJSONLogger(
			// 日志等级
			logger.WithDebugLevel(),
			// 写出的文件
			logger.WithFileRotationP(file),
			// 不在控制台打印
			logger.WithDisableConsole(),
			// 时间格式化
			logger.WithTimeLayout("2006-01-02 15:04:05"),
		)

		// 定时每天凌晨00:00:00初始化日志，让写入的文件名称得以更新
		c.AddFunc("0 0 0 * * *", func() {
			file = "./log/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".log"
			logger.NewJSONLogger(
				// 日志等级
				logger.WithDebugLevel(),
				// 写出的文件
				logger.WithFileRotationP(file),
				// 不在控制台打印
				logger.WithDisableConsole(),
				// 时间格式化
				logger.WithTimeLayout("2006-01-02 15:04:05"),
			)
		})
	}

	if err != nil {
		panic(err)
	}
}
