package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var interviewPointHandle = handle.NewInterviewPointHandle()

// 注册路由
func InterviewPointRouterRegister(r *gin.RouterGroup) {
	interviewPointGroup := r.Group("interviewPoint")
	{
		interviewPointGroup.POST("", interviewPointHandle.Add)
		interviewPointGroup.DELETE("", interviewPointHandle.Delete)
		interviewPointGroup.DELETE("/delete-batch", interviewPointHandle.DeleteByIds)
		interviewPointGroup.GET("/:id", interviewPointHandle.Find)
		interviewPointGroup.GET("/list", interviewPointHandle.List)
		interviewPointGroup.PUT("", interviewPointHandle.Update)
	}
}
