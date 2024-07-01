package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var sysCasbinHandle = handle.NewSysCasbinHandle()

// 注册路由
func SysCasbinRouterRegister(r *gin.RouterGroup) {
	sysCasbinGroup := r.Group("sysCasbin")
	{
		sysCasbinGroup.POST("", sysCasbinHandle.Add)
		sysCasbinGroup.DELETE("", sysCasbinHandle.Delete)
		sysCasbinGroup.DELETE("/delete-batch", sysCasbinHandle.DeleteByIds)
		sysCasbinGroup.GET("/:id", sysCasbinHandle.Find)
		sysCasbinGroup.GET("/list", sysCasbinHandle.List)
		sysCasbinGroup.PUT("", sysCasbinHandle.Update)
	}
}
