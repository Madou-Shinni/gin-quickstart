package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/gin-gonic/gin"
)

func FileRouterRegister(r *gin.Engine) {
	fileGroup := r.Group("file")
	fileHandle := handle.NewFileHandle(service.NewFileService(data.NewFileRepo()))
	{
		fileGroup.POST("", fileHandle.Add)
		fileGroup.POST("/upload-chunk", fileHandle.UploadChunk)
		fileGroup.GET("/chunkid", fileHandle.Chunkid)
		fileGroup.GET("/merge-chunk", fileHandle.MergeChunk)
		fileGroup.DELETE("", fileHandle.Delete)
		fileGroup.DELETE("/delete-batch", fileHandle.DeleteByIds)
		fileGroup.GET("", fileHandle.Find)
		fileGroup.GET("/list", fileHandle.List)
		fileGroup.PUT("", fileHandle.Update)
	}
}
