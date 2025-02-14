package global

import (
	"context"
	"github.com/Madou-Shinni/gin-quickstart/pkg/sms"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/message_queue"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 全局变量
var (
	DB       *Data
	Rdb      *redis.Client
	Producer *message_queue.AsynqClient
	SMS      sms.ISms
)

type Data struct {
	db *gorm.DB
}

func NewData(db *gorm.DB) *Data {
	return &Data{db: db}
}

type contextTxKey struct{}

// Tx gorm Transaction
func (d *Data) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	if tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB); ok {
		// 嵌套事务处理
		return tx.Transaction(func(tx *gorm.DB) error {
			return fn(ctx)
		})
	}
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) WithContext(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db.WithContext(ctx)
}
