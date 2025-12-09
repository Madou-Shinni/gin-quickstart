package gorm_plugin

import (
	"context"
	"time"

	"github.com/Madou-Shinni/go-logger"
	"go.uber.org/zap"
	gormlog "gorm.io/gorm/logger"
)

// GormLogger 实现gorm.io/gorm/logger.Interface接口
var _ gormlog.Interface = (*GormLogger)(nil)

type GormLogger struct {
	level gormlog.LogLevel
}

// NewGormLogger 创建GORM日志适配器
func NewGormLogger() *GormLogger {
	return &GormLogger{
		level: gormlog.Info,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlog.LogLevel) gormlog.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlog.Info {
		logger.Info(msg, zap.Any("data", data))
	}
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlog.Warn {
		logger.Warn(msg, zap.Any("data", data))
	}
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlog.Error {
		logger.Error(msg, zap.Any("data", data))
	}
}

// Trace 记录SQL执行跟踪日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlog.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if logID, ok := ctx.Value("log-id").(string); ok {
		logFields = append(logFields, zap.String("log-id", logID))
	}

	switch {
	case err != nil && l.level >= gormlog.Error:
		logger.Error("SQL Error", append(logFields, zap.Error(err))...)
	case elapsed > 200*time.Millisecond && l.level >= gormlog.Warn:
		logger.Warn("SQL Slow Query", logFields...)
	case l.level >= gormlog.Info:
		logger.Info("SQL Query", logFields...)
	}
}
