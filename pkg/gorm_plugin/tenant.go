package gorm_plugin

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

// 定义上下文中的租户 ID 键
type tenantIDKey struct{}

// 插件结构
type TenantPlugin struct{}

// 定义模型时，包含租户 ID
type Model struct {
	ID       uint   `gorm:"primarykey"`
	TenantID string `gorm:"index"`
}

// 示例表结构
type User struct {
	Model
	Name string
}

// 初始化租户插件
func NewTenantPlugin() *TenantPlugin {
	return &TenantPlugin{}
}

// 注册钩子以自动添加租户 ID 条件
func (tp *TenantPlugin) Apply(db *gorm.DB) {
	// 在查询前设置租户 ID 条件
	db.Callback().Query().Before("gorm:query").Register("tenant:query", func(db *gorm.DB) {
		if tenantID, ok := db.Statement.Context.Value(tenantIDKey{}).(string); ok && db.Statement.Schema != nil {
			tenantField := db.Statement.Schema.LookUpField("TenantID")
			if tenantField != nil {
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{
							Column: clause.Column{Table: db.Statement.Table, Name: "tenant_id"},
							Value:  tenantID,
						},
					},
				})
			}
		}
	})

	// 在创建前设置租户 ID
	db.Callback().Create().Before("gorm:create").Register("tenant:create", func(db *gorm.DB) {
		if tenantID, ok := db.Statement.Context.Value(tenantIDKey{}).(string); ok && db.Statement.Schema != nil {
			tenantField := db.Statement.Schema.LookUpField("TenantID")
			if tenantField != nil {
				db.Statement.SetColumn("TenantID", tenantID)
			}
		}
	})

	// 在修改前设置租户 ID
	db.Callback().Update().Before("gorm:update").Register("tenant:update", func(db *gorm.DB) {
		if tenantID, ok := db.Statement.Context.Value(tenantIDKey{}).(string); ok && db.Statement.Schema != nil {
			tenantField := db.Statement.Schema.LookUpField("TenantID")
			if tenantField != nil {
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{
							Column: clause.Column{Table: db.Statement.Table, Name: "tenant_id"},
							Value:  tenantID,
						},
					},
				})
			}
		}
	})

	// 在删除前设置租户 ID
	db.Callback().Delete().Before("gorm:delete").Register("tenant:delete", func(db *gorm.DB) {
		if tenantID, ok := db.Statement.Context.Value(tenantIDKey{}).(string); ok && db.Statement.Schema != nil {
			tenantField := db.Statement.Schema.LookUpField("TenantID")
			if tenantField != nil {
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{
							Column: clause.Column{Table: db.Statement.Table, Name: "tenant_id"},
							Value:  tenantID,
						},
					},
				})
			}
		}
	})
}

// 中间件，提取租户 ID 并存储在 Gin 的上下文中
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//tenantID := c.GetHeader("X-Tenant-ID")
		tenantID := "1"
		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "X-Tenant-ID header is required"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), tenantIDKey{}, tenantID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
