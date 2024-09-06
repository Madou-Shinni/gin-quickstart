package main

import (
	"context"
	_ "github.com/Madou-Shinni/gin-quickstart/initialize"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/go-openapi/loads"
	"gorm.io/gorm/clause"
	"log"
)

var (
	defaultDocsPath = "docs/swagger.yaml"
)

func main() {
	// 同步api到库
	syncApi()
}

func syncApi() {
	// Load the swagger document
	doc, err := loads.Spec(defaultDocsPath)
	if err != nil {
		log.Fatalf("Failed to load spec: %v", err)
	}

	// Get the Swagger object
	swagger := doc.Spec()

	var slice []*domain.SysApi
	// Print all paths
	for path, pathItem := range swagger.Paths.Paths {
		//fmt.Printf("Path: %s\n", path)
		item := domain.SysApi{Path: path}
		// Print the methods (GET, POST, etc.) and their details
		if pathItem.Get != nil {
			//fmt.Printf("  Method: GET\n    Summary: %s\n", pathItem.Get.Summary)
			item.Method = "GET"
			item.Name = pathItem.Get.Summary
		}
		if pathItem.Post != nil {
			//fmt.Printf("  Method: POST\n    Summary: %s\n", pathItem.Post.Summary)
			item.Method = "POST"
			item.Name = pathItem.Post.Summary
		}
		if pathItem.Put != nil {
			//fmt.Printf("  Method: PUT\n    Summary: %s\n", pathItem.Put.Summary)
			item.Method = "PUT"
			item.Name = pathItem.Put.Summary
		}
		if pathItem.Delete != nil {
			//fmt.Printf("  Method: DELETE\n    Summary: %s\n", pathItem.Delete.Summary)
			item.Method = "DELETE"
			item.Name = pathItem.Delete.Summary
		}

		slice = append(slice, &item)
	}

	err = global.DB.WithContext(context.Background()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "method"}, {Name: "path"}},
		DoUpdates: clause.AssignmentColumns([]string{"method", "path", "name"}),
	}).Create(&slice).Error
	if err != nil {
		log.Fatalf("Failed to save docs: %v", err)
		return
	}

	log.Println("Saved docs success")
}
