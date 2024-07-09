package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/Madou-Shinni/gin-quickstart/middleware"
	"github.com/gin-gonic/gin"
)

var sysUserHandle = handle.NewSysUserHandle()

// 注册路由
func SysUserRouterRegister(r *gin.RouterGroup) {
	sysUserGroup := r.Group("sysUser", middleware.JwtAuth(), middleware.CasbinHandler())
	{
		sysUserGroup.POST("", sysUserHandle.Add)
		sysUserGroup.DELETE("", sysUserHandle.Delete)
		sysUserGroup.DELETE("/delete-batch", sysUserHandle.DeleteByIds)
		sysUserGroup.GET("/:id", sysUserHandle.Find)
		sysUserGroup.GET("/list", sysUserHandle.List)
		sysUserGroup.GET("/info", sysUserHandle.Info)
		sysUserGroup.PUT("", sysUserHandle.Update)
	}

	sysUserGroupNoAuth := r.Group("sysUser")
	{
		sysUserGroupNoAuth.POST("/login", sysUserHandle.Login)
	}
}
