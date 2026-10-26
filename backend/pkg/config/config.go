package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the process configuration.
type Config struct {
	HTTPAddr    string `yaml:"http_addr"`
	DatabaseURL string `yaml:"database_url"`
	StaticDir   string `yaml:"static_dir"`
}

// Load reads YAML from CONFIG_PATH when that variable is set, then applies
// environment variables. Compose sets the environment.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr: ":8080",
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

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	return cfg, nil
}

func override(dst *string, key string) {
	if v := os.Getenv(key); v != "" {
		*dst = v
	}
}
