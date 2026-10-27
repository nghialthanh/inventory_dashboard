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
	"inventory_dashboard/internal/common/snowflake"
	"inventory_dashboard/internal/logic/catalog"
	"inventory_dashboard/internal/logic/health"
	"inventory_dashboard/internal/logic/vehicle"
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/pkg/cache"
	"inventory_dashboard/pkg/config"
	"inventory_dashboard/pkg/database"
	"inventory_dashboard/pkg/logger"
)

// App is the process: PostgreSQL, optional Redis, and the HTTP server.
type App struct {
	db     *pgxpool.Pool
	readDB *pgxpool.Pool
	redis  *cache.Redis
}

// Run connects to PostgreSQL (and optional Redis), then serves the API and the screen.
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

	readDB := db
	if cfg.DatabaseReadURL != cfg.DatabaseURL {
		replica, err := retry(ctx, func() (*pgxpool.Pool, error) {
			return database.Open(ctx, cfg.DatabaseReadURL)
		})
		if err != nil {
			db.Close()
			return fmt.Errorf("database read url: %w", err)
		}
		readDB = replica
	}

	application := &App{db: db, readDB: readDB}
	defer application.Close()

	var catalogCache catalog.Cache
	if cfg.RedisURL != "" {
		rdb, err := retry(ctx, func() (*cache.Redis, error) {
			r, err := cache.Open(cfg.RedisURL)
			if err != nil {
				return nil, err
			}
			if err := r.Ping(ctx); err != nil {
				_ = r.Close()
				return nil, err
			}
			return r, nil
		})
		if err != nil {
			return fmt.Errorf("redis: %w", err)
		}
		application.redis = rdb
		catalogCache = rdb
		logger.L().Info("catalog cache enabled",
			zap.Duration("ttl", cfg.CatalogCacheTTL),
		)
	}

	ids, err := snowflake.New(cfg.WorkerID, nil)
	if err != nil {
		return err
	}
	catalogRepo := repositories.NewCatalog(db, readDB)
	e := routes.New(routes.Config{
		StaticDir: cfg.StaticDir,
		Health:    health.New(),
		Catalog:   catalog.New(catalogRepo, catalogCache, cfg.CatalogCacheTTL),
		Vehicles:  vehicle.New(repositories.NewVehicles(db), catalogRepo, ids, nil),
	})
	e.Server.ReadHeaderTimeout = 5 * time.Second
	logger.L().Info("inventory api listening",
		zap.String("addr", cfg.HTTPAddr),
		zap.Duration("shutdown_timeout", cfg.ShutdownTimeout),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Start(cfg.HTTPAddr)
	}()

	select {
	case <-ctx.Done():
		logger.L().Info("graceful shutdown started")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			return err
		}
		logger.L().Info("graceful shutdown complete")
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

// Close releases database pools and Redis.
func (a *App) Close() {
	if a.redis != nil {
		_ = a.redis.Close()
	}
	if a.readDB != nil && a.readDB != a.db {
		a.readDB.Close()
	}
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
