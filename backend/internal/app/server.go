package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"inventory_dashboard/internal/app/routes"
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/internal/services"
	"inventory_dashboard/pkg/config"
	"inventory_dashboard/pkg/database"
	"inventory_dashboard/pkg/storage"
)

// App is the process: PostgreSQL, private MinIO, and the HTTP server.
type App struct {
	db    *pgxpool.Pool
	media *storage.Client
}

// Run connects to PostgreSQL and MinIO, then serves the API and the screen.
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
	media, err := storage.Open(storage.Settings{
		Endpoint:  cfg.MinIOEndpoint,
		AccessKey: cfg.MinIOAccessKey,
		SecretKey: cfg.MinIOSecretKey,
		Bucket:    cfg.MinIOBucket,
		UseSSL:    cfg.MinIOUseSSL,
	})
	if err != nil {
		db.Close()
		return err
	}
	if err := retryErr(ctx, media.EnsureBucket); err != nil {
		db.Close()
		return err
	}

	application := &App{db: db, media: media}
	defer application.Close()

	svc := services.New(repositories.NewCatalog(db))
	e := routes.New(routes.Config{StaticDir: cfg.StaticDir, Services: svc})
	e.Server.ReadHeaderTimeout = 5 * time.Second
	slog.Info("inventory api listening", "addr", cfg.HTTPAddr, "bucket", application.media.Bucket())

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

func retryErr(ctx context.Context, op func(context.Context) error) error {
	_, err := retry(ctx, func() (struct{}, error) {
		return struct{}{}, op(ctx)
	})
	return err
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
