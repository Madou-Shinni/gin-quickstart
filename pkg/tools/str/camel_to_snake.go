package str

import (
	"regexp"
	"strings"
)

// CamelToSnake 驼峰转下划线
func CamelToSnake(s string) string {
	re := regexp.MustCompile("[A-Z]")
	snake := re.ReplaceAllStringFunc(s, func(match string) string {
		return "_" + strings.ToLower(match)
	})
	return strings.TrimPrefix(strings.ToLower(snake), "_")
}
