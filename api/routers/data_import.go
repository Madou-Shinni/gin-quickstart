package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var dataImportHandle = handle.NewDataImportHandle()

// 注册路由
func DataImportRouterRegister(r *gin.RouterGroup) {
	dataImportGroup := r.Group("dataImport")
	{
		dataImportGroup.POST("", dataImportHandle.Add)
		dataImportGroup.DELETE("", dataImportHandle.Delete)
		dataImportGroup.DELETE("/delete-batch", dataImportHandle.DeleteByIds)
		dataImportGroup.GET("/:id", dataImportHandle.Find)
		dataImportGroup.GET("/list", dataImportHandle.List)
		dataImportGroup.GET("/template", dataImportHandle.Template)
		dataImportGroup.POST("/import", dataImportHandle.Import)
		dataImportGroup.PUT("", dataImportHandle.Update)
	}
}
