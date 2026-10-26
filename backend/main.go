package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"inventory_dashboard/internal/app"
	"inventory_dashboard/pkg/logger"
)

func main() {
	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := app.Run(ctx); err != nil {
		logger.L().Error("inventory api stopped", zap.Error(err))
		os.Exit(1)
	}
}
