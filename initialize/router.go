package initialize

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/api/routers"
	_ "github.com/Madou-Shinni/gin-quickstart/docs"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/middleware"
	"github.com/Madou-Shinni/go-logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// RunServer 启动服务
func RunServer() {
	// 初始化引擎
	if conf.Conf.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 设置 swagger 访问路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 跨域
	r.Use(cors.Default())

	// 队列监控
	asynqConfig := conf.Conf.AsynqConfig
	amconfig := asynqConfig.Monitor
	if amconfig.Enable {
		h := asynqmon.New(asynqmon.Options{
			RootPath:     amconfig.RootPath, // RootPath specifies the root for asynqmon app
			RedisConnOpt: asynq.RedisClientOpt{Addr: asynqConfig.Addr, Password: asynqConfig.Password, DB: asynqConfig.DB},
			ReadOnly:     amconfig.Readonly,
		})
		r.Any(h.RootPath()+"/*any", middleware.UnameAndPwdAuth(), gin.WrapH(h))
	}

	// 缓存
	//r.Use(middleware.Cache(cache.NewRdbCache(global.Rdb)))

	// 设置路由组
	public := r.Group("")
	private := r.Group("", middleware.JwtAuth(), middleware.CasbinHandler())

	// 注册路由
	// 热更新日志级别 debug info warn error
	r.PUT("/logs-lvl", gin.WrapH(logger.ChangeLevelHandlerFunc()))
	routers.DemoRouterRegister(public)
	routers.FileRouterRegister(r)
	routers.SystemRouterRegister(public)
	routers.SysUserRouterRegister(public)
	routers.SysRoleRouterRegister(public)
	routers.SysCasbinRouterRegister(public)
	routers.SysApiRouterRegister(private)
	routers.SysMenuRouterRegister(private)
	routers.DataImportRouterRegister(public)
	routers.NoPageRouterRegister(private)

	log.Printf("[GIN-QuickStart] 接口文档地址：http://localhost:%v/swagger/index.html\n", conf.Conf.ServerPort)

	r.Run(fmt.Sprintf("0.0.0.0:%v", conf.Conf.ServerPort))
}
