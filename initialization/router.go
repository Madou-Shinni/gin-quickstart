package initialization

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/api/routers"
	_ "github.com/Madou-Shinni/gin-quickstart/docs"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/middleware"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/cache"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 启动服务
func RunServer() {
	// 初始化引擎
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 设置 swagger 访问路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 跨域
	r.Use(cors.Default())

	// 缓存
	r.Use(middleware.Cache(cache.NewRdbCache(global.Rdb)))

	// 注册路由
	routers.DemoRouterRegister(r)
	routers.FileRouterRegister(r)

	fmt.Printf("[GIN-QuickStart] 接口文档地址：http://localhost:%v/swagger/index.html\n", conf.Conf.ServerPort)

	r.Run(fmt.Sprintf("0.0.0.0:%v", conf.Conf.ServerPort))
}
