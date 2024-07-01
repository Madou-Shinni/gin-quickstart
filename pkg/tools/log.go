package tools

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

// CustomTimeBasedLogger 实现按天分割的日志记录器
type CustomTimeBasedLogger struct {
	logger     *zap.Logger
	rotateDate time.Time
}

func NewCustomTimeBasedLogger() *CustomTimeBasedLogger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   getLogFilename(),
		MaxSize:    100, // 保留原有的大小分割机制
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}
	lumberjackLogger.Rotate()

	w := zapcore.AddSync(lumberjackLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		zap.InfoLevel,
	)

	logger := zap.New(core)

	return &CustomTimeBasedLogger{
		logger:     logger,
		rotateDate: time.Now().Truncate(24 * time.Hour),
	}
}

func getLogFilename() string {
	return fmt.Sprintf("./logs/app-%s.log", time.Now().Format(time.DateOnly))
}

func (ctl *CustomTimeBasedLogger) rotateIfNeeded() {
	currentDate := time.Now().Truncate(24 * time.Hour)
	if currentDate.After(ctl.rotateDate) {
		ctl.rotateDate = currentDate
		ctl.logger.Sync() // flushes buffer, if any

		// Reinitialize logger with new filename
		ctl.logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   getLogFilename(),
				MaxSize:    100,
				MaxBackups: 7,
				MaxAge:     30,
				Compress:   true,
			}),
			zap.InfoLevel,
		))
	}
}

func (ctl *CustomTimeBasedLogger) Info(msg string, fields ...zapcore.Field) {
	ctl.rotateIfNeeded()
	ctl.logger.Info(msg, fields...)
}

func main() {
	logger := NewCustomTimeBasedLogger()

	for i := 0; i < 1000; i++ {
		logger.Info("Logging with zap and lumberjack",
			zap.Int("iteration", i),
		)
		time.Sleep(1 * time.Second)
	}
}
