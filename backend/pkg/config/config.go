package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the process configuration. MinIO credentials stay on the Compose
// network; the screen never receives them.
type Config struct {
	HTTPAddr       string `yaml:"http_addr"`
	DatabaseURL    string `yaml:"database_url"`
	StaticDir      string `yaml:"static_dir"`
	MinIOEndpoint  string `yaml:"minio_endpoint"`
	MinIOAccessKey string `yaml:"minio_access_key"`
	MinIOSecretKey string `yaml:"minio_secret_key"`
	MinIOBucket    string `yaml:"minio_bucket"`
	MinIOUseSSL    bool   `yaml:"minio_use_ssl"`
}

// Load reads YAML from CONFIG_PATH when that variable is set, then applies
// environment variables. Compose sets the environment and does not publish MinIO.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:    ":8080",
		MinIOBucket: "vehicle-images",
	}

	if path := os.Getenv("CONFIG_PATH"); path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
		if err := yaml.Unmarshal(raw, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
	}

	override(&cfg.HTTPAddr, "HTTP_ADDR")
	override(&cfg.DatabaseURL, "DATABASE_URL")
	override(&cfg.StaticDir, "STATIC_DIR")
	override(&cfg.MinIOEndpoint, "MINIO_ENDPOINT")
	override(&cfg.MinIOAccessKey, "MINIO_ACCESS_KEY")
	override(&cfg.MinIOSecretKey, "MINIO_SECRET_KEY")
	override(&cfg.MinIOBucket, "MINIO_BUCKET")
	if v := os.Getenv("MINIO_USE_SSL"); v != "" {
		cfg.MinIOUseSSL = v == "true"
	}

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if cfg.MinIOBucket == "" {
		cfg.MinIOBucket = "vehicle-images"
	}
	if cfg.DatabaseURL == "" || cfg.MinIOEndpoint == "" || cfg.MinIOAccessKey == "" || cfg.MinIOSecretKey == "" {
		return Config{}, fmt.Errorf("DATABASE_URL, MINIO_ENDPOINT, MINIO_ACCESS_KEY, and MINIO_SECRET_KEY are required")
	}
	return cfg, nil
}

func override(dst *string, key string) {
	if v := os.Getenv(key); v != "" {
		*dst = v
	}
}
