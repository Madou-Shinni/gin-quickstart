package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var sysMenuHandle = handle.NewSysMenuHandle()

// 注册路由
func SysMenuRouterRegister(r *gin.RouterGroup) {
	sysMenuGroup := r.Group("sysMenu")
	{
		sysMenuGroup.POST("", sysMenuHandle.Add)
		sysMenuGroup.DELETE("", sysMenuHandle.Delete)
		sysMenuGroup.DELETE("/delete-batch", sysMenuHandle.DeleteByIds)
		sysMenuGroup.GET("/:id", sysMenuHandle.Find)
		sysMenuGroup.GET("/list", sysMenuHandle.List)
		sysMenuGroup.GET("/role-list", sysMenuHandle.RoleList)
		sysMenuGroup.PUT("", sysMenuHandle.Update)
		sysMenuGroup.PUT("/role-list", sysMenuHandle.SetRoleList)
	}
}
