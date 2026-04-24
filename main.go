package main

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// init logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Fatal to init logger: %v", err)
	}
	logger.Info("Starting server...")

	// init config from .env or environment
	if err := godotenv.Load(); err != nil {
		logger.Info("Not found .env file. Using existing environment")
	}
	cfg, err := load_config()
	if err != nil {
		logger.Error("Failed to load configuration")
	}

	db, err := init_db(cfg, logger)
	if err != nil {
		logger.Error("Failed to connect to database")
	}

	if err := auto_migrate(db); err != nil {
		logger.Error("Automigration failed")
	}

	// init crudl route

	// init cors

	// init docs (swagger)

}
