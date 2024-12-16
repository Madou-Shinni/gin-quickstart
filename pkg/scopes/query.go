package scopes

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
	"gorm.io/gorm"
)

// Paginate 分页
// 默认加载 pagelimit.OffsetLimit 的分页
func Paginate(page request.PageSearch) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 翻页
		if page.NoPage {
			return db
		}
		offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)
		return db.Offset(offset).Limit(limit)
	}
}
