package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var systemHandle = handle.NewSystemHandle()

// 注册路由
func SystemRouterRegister(r *gin.RouterGroup) {
	systemGroup := r.Group("system")
	{
		systemGroup.POST("init", systemHandle.Init)
	}
}
