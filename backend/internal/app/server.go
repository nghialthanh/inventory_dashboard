package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"inventory_dashboard/internal/app/routes"
	"inventory_dashboard/internal/logic/catalog"
	"inventory_dashboard/internal/logic/health"
	"inventory_dashboard/internal/logic/vehicle"
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/pkg/config"
	"inventory_dashboard/pkg/database"
	"inventory_dashboard/pkg/logger"
)

// App is the process: PostgreSQL and the HTTP server.
type App struct {
	db *pgxpool.Pool
}

// Run connects to PostgreSQL, then serves the API and the screen.
func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.StaticDir != "" {
		if _, err := os.Stat(filepath.Join(cfg.StaticDir, "index.html")); err != nil {
			return fmt.Errorf("static dir: %w", err)
		}
	}

	db, err := retry(ctx, func() (*pgxpool.Pool, error) {
		return database.Open(ctx, cfg.DatabaseURL)
	})
	if err != nil {
		return err
	}
	if err := database.Migrate(ctx, db); err != nil {
		db.Close()
		return err
	}

	application := &App{db: db}
	defer application.Close()

	catalogRepo := repositories.NewCatalog(db)
	e := routes.New(routes.Config{
		StaticDir: cfg.StaticDir,
		Health:    health.New(),
		Catalog:   catalog.New(catalogRepo),
		Vehicles:  vehicle.New(repositories.NewVehicles(db), catalogRepo, nil),
	})
	e.Server.ReadHeaderTimeout = 5 * time.Second
	logger.L().Info("inventory api listening",
		zap.String("addr", cfg.HTTPAddr),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Start(cfg.HTTPAddr)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

// Close releases the database pool.
func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}

func retry[T any](ctx context.Context, op func() (T, error)) (T, error) {
	var zero T
	var last error
	for attempt := 1; attempt <= 30; attempt++ {
		value, err := op()
		if err == nil {
			return value, nil
		}
		last = err
		if !sleep(ctx, time.Second) {
			return zero, last
		}
	}
	return zero, last
}

func sleep(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
