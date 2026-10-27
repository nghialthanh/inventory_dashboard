package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the process configuration.
type Config struct {
	HTTPAddr        string
	DatabaseURL     string
	DatabaseReadURL string
	RedisURL        string
	CatalogCacheTTL time.Duration
	ShutdownTimeout time.Duration
	StaticDir       string
	WorkerID        int64
}

type fileConfig struct {
	HTTPAddr        string `yaml:"http_addr"`
	DatabaseURL     string `yaml:"database_url"`
	DatabaseReadURL string `yaml:"database_read_url"`
	RedisURL        string `yaml:"redis_url"`
	CatalogCacheTTL string `yaml:"catalog_cache_ttl"`
	ShutdownTimeout string `yaml:"shutdown_timeout"`
	StaticDir       string `yaml:"static_dir"`
	WorkerID        int64  `yaml:"worker_id"`
}

// Load reads YAML from CONFIG_PATH when that variable is set, then applies
// environment variables. Compose sets the environment.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:        ":8080",
		WorkerID:        0,
		CatalogCacheTTL: time.Hour,
		ShutdownTimeout: 15 * time.Second,
	}

	if path := os.Getenv("CONFIG_PATH"); path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
		var file fileConfig
		if err := yaml.Unmarshal(raw, &file); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
		if file.HTTPAddr != "" {
			cfg.HTTPAddr = file.HTTPAddr
		}
		if file.DatabaseURL != "" {
			cfg.DatabaseURL = file.DatabaseURL
		}
		if file.DatabaseReadURL != "" {
			cfg.DatabaseReadURL = file.DatabaseReadURL
		}
		if file.RedisURL != "" {
			cfg.RedisURL = file.RedisURL
		}
		if file.StaticDir != "" {
			cfg.StaticDir = file.StaticDir
		}
		cfg.WorkerID = file.WorkerID
		if file.CatalogCacheTTL != "" {
			d, err := time.ParseDuration(file.CatalogCacheTTL)
			if err != nil {
				return Config{}, fmt.Errorf("catalog_cache_ttl: %w", err)
			}
			cfg.CatalogCacheTTL = d
		}
		if file.ShutdownTimeout != "" {
			d, err := time.ParseDuration(file.ShutdownTimeout)
			if err != nil {
				return Config{}, fmt.Errorf("shutdown_timeout: %w", err)
			}
			cfg.ShutdownTimeout = d
		}
	}

	override(&cfg.HTTPAddr, "HTTP_ADDR")
	override(&cfg.DatabaseURL, "DATABASE_URL")
	override(&cfg.DatabaseReadURL, "DATABASE_READ_URL")
	override(&cfg.RedisURL, "REDIS_URL")
	override(&cfg.StaticDir, "STATIC_DIR")
	if raw := os.Getenv("WORKER_ID"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return Config{}, fmt.Errorf("WORKER_ID: %w", err)
		}
		cfg.WorkerID = id
	}
	if err := overrideDuration(&cfg.CatalogCacheTTL, "CATALOG_CACHE_TTL"); err != nil {
		return Config{}, err
	}
	if err := overrideDuration(&cfg.ShutdownTimeout, "SHUTDOWN_TIMEOUT"); err != nil {
		return Config{}, err
	}

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.DatabaseReadURL == "" {
		cfg.DatabaseReadURL = cfg.DatabaseURL
	}
	if cfg.CatalogCacheTTL <= 0 {
		cfg.CatalogCacheTTL = time.Hour
	}
	if cfg.ShutdownTimeout <= 0 {
		cfg.ShutdownTimeout = 15 * time.Second
	}
	return cfg, nil
}

func override(dst *string, key string) {
	if v := os.Getenv(key); v != "" {
		*dst = v
	}
}

func overrideDuration(dst *time.Duration, key string) error {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return fmt.Errorf("%s: %w", key, err)
	}
	*dst = d
	return nil
}
