package scopes

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// MatchStringSliceScope returns a GORM scope to filter JSON tags.
// - vals: 字符串数组
// - matchAll: 是否要求全部匹配 (true 为 AND，false 为 OR)
func MatchStringSliceScope(col string, vals []string, matchAll bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(vals) == 0 {
			return db
		}
		conditions := make([]string, 0, len(vals))
		args := make([]interface{}, 0, len(vals))

		for _, tag := range vals {
			conditions = append(conditions, fmt.Sprintf("JSON_CONTAINS(%s, ?)", col))
			args = append(args, fmt.Sprintf(`["%s"]`, tag))
		}

		join := " OR "
		if matchAll {
			join = " AND "
		}
		return db.Where(strings.Join(conditions, join), args...)
	}
}
