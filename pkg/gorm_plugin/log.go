package gorm_plugin

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

// 定义上下文中的租户 ID 键
type logIDKey struct{}

// 插件结构
type logPlugin struct{}

// 初始化租户插件
func NewLogPlugin() *logPlugin {
	return &logPlugin{}
}

// 注册钩子以自动添加租户 ID 条件
func (tp *logPlugin) Apply(db *gorm.DB) {
	// 在查询前设置租户 ID 条件
	db.Callback().Query().After("gorm:query").Register("log:query", func(db *gorm.DB) {
		if logID, ok := db.Statement.Context.Value(logIDKey{}).(string); ok && db.Statement.Schema != nil {
			explain := db.Statement.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
			log.Printf("log id %s sql %s", logID, explain)
		}
	})
}

// 中间件，提取租户 ID 并存储在 Gin 的上下文中
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logID := c.GetHeader("X-log-ID")
		if logID != "" {
			ctx := context.WithValue(c.Request.Context(), logIDKey{}, logID)
			ctx = context.WithValue(ctx, "log-id", logID)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}
	}
}
