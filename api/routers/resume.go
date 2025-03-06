package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var resumeHandle = handle.NewResumeHandle()

// 注册路由
func ResumeRouterRegister(r *gin.RouterGroup) {
	resumeGroup := r.Group("resume")
	{
		resumeGroup.POST("", resumeHandle.Add)
		resumeGroup.DELETE("", resumeHandle.Delete)
		resumeGroup.DELETE("/delete-batch", resumeHandle.DeleteByIds)
		resumeGroup.GET("/:id", resumeHandle.Find)
		resumeGroup.GET("/list", resumeHandle.List)
		resumeGroup.PUT("", resumeHandle.Update)
	}
}
