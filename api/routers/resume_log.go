package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var resumeLogHandle = handle.NewResumeLogHandle()

// 注册路由
func ResumeLogRouterRegister(r *gin.RouterGroup) {
	resumeLogGroup := r.Group("resumeLog")
	{
		resumeLogGroup.POST("", resumeLogHandle.Add)
		resumeLogGroup.POST("/push", resumeLogHandle.Push)
		resumeLogGroup.DELETE("", resumeLogHandle.Delete)
		resumeLogGroup.DELETE("/delete-batch", resumeLogHandle.DeleteByIds)
		resumeLogGroup.GET("/:id", resumeLogHandle.Find)
		resumeLogGroup.GET("/list", resumeLogHandle.List)
		resumeLogGroup.PUT("", resumeLogHandle.Update)
	}
}
