package gorm_plugin

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestTenantPlugin_Apply(t *testing.T) {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移数据库
	db.AutoMigrate(&User{})

	// 应用租户插件
	tenantPlugin := NewTenantPlugin()
	tenantPlugin.Apply(db)

	// 设置 Gin
	r := gin.Default()

	// 使用租户中间件
	r.Use(TenantMiddleware())

	// 路由处理
	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 使用包含租户 ID 的上下文
		ctx := c.Request.Context()
		if err := db.Debug().WithContext(ctx).Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.GET("/users", func(c *gin.Context) {
		var users []User

		// 使用包含租户 ID 的上下文
		ctx := c.Request.Context()
		if err := db.Debug().WithContext(ctx).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
