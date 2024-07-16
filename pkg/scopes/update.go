package scopes

import "gorm.io/gorm"

// UpdatesAllOmit 更新所有的字段包括零值
// 默认排除 id created_at
// 可以额外添加排除字段
func UpdatesAllOmit(fields ...string) func(db *gorm.DB) *gorm.DB {
	defaultOmit := []string{"id", "created_at"}
	return func(db *gorm.DB) *gorm.DB {
		fields = append(fields, defaultOmit...)
		return db.Select("*").Omit(fields...)
	}
}
