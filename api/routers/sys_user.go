package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var sysUserHandle = handle.NewSysUserHandle()

// 注册路由
func SysUserRouterRegister(r *gin.RouterGroup) {
	sysUserGroup := r.Group("sysUser")
	{
		sysUserGroup.POST("", sysUserHandle.Add)
		sysUserGroup.DELETE("", sysUserHandle.Delete)
		sysUserGroup.DELETE("/delete-batch", sysUserHandle.DeleteByIds)
		sysUserGroup.GET("/:id", sysUserHandle.Find)
		sysUserGroup.GET("/list", sysUserHandle.List)
		sysUserGroup.PUT("", sysUserHandle.Update)
	}
}
