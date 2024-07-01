package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var sysApiHandle = handle.NewSysApiHandle()

// 注册路由
func SysApiRouterRegister(r *gin.RouterGroup) {
	sysApiGroup := r.Group("sysApi")
	{
		sysApiGroup.POST("", sysApiHandle.Add)
		sysApiGroup.DELETE("", sysApiHandle.Delete)
		sysApiGroup.DELETE("/delete-batch", sysApiHandle.DeleteByIds)
		sysApiGroup.GET("/:id", sysApiHandle.Find)
		sysApiGroup.GET("/list", sysApiHandle.List)
		sysApiGroup.PUT("", sysApiHandle.Update)
	}
}
