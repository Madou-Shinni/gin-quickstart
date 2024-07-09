package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var sysRoleHandle = handle.NewSysRoleHandle()

// 注册路由
func SysRoleRouterRegister(r *gin.RouterGroup) {
	sysRoleGroup := r.Group("sysRole")
	{
		sysRoleGroup.POST("", sysRoleHandle.Add)
		sysRoleGroup.DELETE("", sysRoleHandle.Delete)
		sysRoleGroup.DELETE("/delete-batch", sysRoleHandle.DeleteByIds)
		sysRoleGroup.GET("/:id", sysRoleHandle.Find)
		sysRoleGroup.GET("/list", sysRoleHandle.List)
		sysRoleGroup.PUT("", sysRoleHandle.Update)
		sysRoleGroup.PUT("/user-list", sysRoleHandle.SetUserRoleList)
	}
}
