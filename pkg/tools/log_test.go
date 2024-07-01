package tools

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestRotation(t *testing.T) {
	logger := NewCustomTimeBasedLogger()

	for i := 0; i < 100; i++ {
		logger.Info("Logging with zap and lumberjack",
			zap.Int("iteration", i),
		)
		time.Sleep(1 * time.Second)
	}
}
