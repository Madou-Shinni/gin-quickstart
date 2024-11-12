package main

// 使用 _引入依赖项在main函数执行会直接调用init函数
import (
	"github.com/Madou-Shinni/gin-quickstart/initialize"
	"github.com/Madou-Shinni/gin-quickstart/job"
	"github.com/Madou-Shinni/gin-quickstart/task"
	"github.com/Madou-Shinni/go-logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 生成swagger文档
// --parseDependency --parseInternal 识别到内外部依赖
// --parseDepth 解析依赖包的深度默认100 设置深度可以大大加快生成速度(推荐使用这个)
// --output 文件生成目录
//go:generate swag initialize --parseDependency --parseInternal --output ../docs
//go:generate swag initialize --parseDependency --parseInternal --parseDepth 3 --output ../docs
//go:generate swag initialize --output ../docs

func main() {
	// 启动服务(使用goroutine解决服务启动时程序阻塞问题)
	go initialize.RunServer()
	go job.RunConsumer()
	go task.V2Init()
	// 监听信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-signals:
		// 释放资源
		logger.Sync()
		initialize.Close()
		log.Println("[GIN-QuickStart] 程序关闭，释放资源")
		return
	}
}
