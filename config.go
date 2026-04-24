package main

import (
	"os"
)

type Config struct {
	DatabaseURL string
}

func load_config() (*Config, error) {
	get := func(key, default_value string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return default_value
	}
	return &Config{
		DatabaseURL: get("DATABASE_URL", "postgres://user:password@localhost:5432/core_db?sslmode=disable"),
	}, nil
}
