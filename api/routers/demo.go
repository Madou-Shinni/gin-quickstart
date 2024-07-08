package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/Madou-Shinni/gin-quickstart/pkg/gorm_plugin"
	"github.com/gin-gonic/gin"
)

var demoHandle = handle.NewDemoHandle()

// 注册路由
func DemoRouterRegister(r *gin.RouterGroup) {
	demoGroup := r.Group("demo", gorm_plugin.LogMiddleware())
	{
		demoGroup.POST("", demoHandle.Add)
		demoGroup.DELETE("", demoHandle.Delete)
		demoGroup.DELETE("/delete-batch", demoHandle.DeleteByIds)
		demoGroup.GET("/:id", demoHandle.Find)
		demoGroup.GET("/list", demoHandle.List)
		demoGroup.PUT("", demoHandle.Update)
	}
}
