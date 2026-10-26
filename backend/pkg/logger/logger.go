package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"inventory_dashboard/internal/common/trace"
)

var (
	mu     sync.RWMutex
	global *zap.Logger = zap.NewNop()
)

// Init configures the process logger. LOG_DEV=true uses a console encoder.
func Init() error {
	var (
		cfg zap.Config
		err error
	)
	if os.Getenv("LOG_DEV") == "true" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.MessageKey = "msg"
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	built, err := cfg.Build(zap.AddCallerSkip(0))
	if err != nil {
		return err
	}
	mu.Lock()
	if global != nil {
		_ = global.Sync()
	}
	global = built
	mu.Unlock()
	zap.ReplaceGlobals(built)
	return nil
}

// L returns the process logger.
func L() *zap.Logger {
	mu.RLock()
	defer mu.RUnlock()
	return global
}

// FromContext returns a logger with trace_id when present on ctx.
func FromContext(ctx context.Context) *zap.Logger {
	log := L()
	if id := trace.ID(ctx); id != "" {
		return log.With(zap.String("trace_id", id))
	}
	return log
}

// Sync flushes buffered entries. Safe to call when uninitialized.
func Sync() {
	mu.RLock()
	log := global
	mu.RUnlock()
	if log != nil {
		_ = log.Sync()
	}
}
