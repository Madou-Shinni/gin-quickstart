package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/gin-gonic/gin"
)

var fileHandle = handle.NewFileHandle()

func FileRouterRegister(r *gin.Engine) {
	fileGroup := r.Group("file")
	{
		fileGroup.POST("", fileHandle.Add)
		fileGroup.POST("/upload-chunk", fileHandle.UploadChunk)
		fileGroup.POST("/upload", fileHandle.Upload)
		fileGroup.GET("/chunkid", fileHandle.Chunkid)
		fileGroup.GET("/merge-chunk", fileHandle.MergeChunk)
		fileGroup.DELETE("", fileHandle.Delete)
		fileGroup.DELETE("/delete-batch", fileHandle.DeleteByIds)
		fileGroup.GET("/:id", fileHandle.Find)
		fileGroup.GET("/list", fileHandle.List)
		fileGroup.PUT("", fileHandle.Update)
	}
}
