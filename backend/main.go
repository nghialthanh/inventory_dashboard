package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"inventory_dashboard/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := app.Run(ctx); err != nil {
		slog.Error("inventory api stopped", "err", err)
		os.Exit(1)
	}
}
