package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

func NoPageRouterRegister(r *gin.RouterGroup) {
	demoGroup := r.Group("noPage")
	{
		demoGroup.GET("/:table", handle.NoPage)
	}
}
