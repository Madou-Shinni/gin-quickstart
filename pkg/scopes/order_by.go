package scopes

import "gorm.io/gorm"

// OrderBy 根据字段排序
func OrderBy(str string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(str)
	}
}
