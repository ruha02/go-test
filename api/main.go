package main

import (
	"log"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	// init logger
	logger, err := zap.NewDevelopment()

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
		return
	}
	logger.Info("Environment", zap.String("DATABASE_URL", cfg.DatabaseURL))

	db, err := init_db(cfg, logger)
	if err != nil {
		logger.Error("Failed to connect to database")
		return
	}

	if err := auto_migrate(db); err != nil {
		logger.Error("Automigration failed")
	}

	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	createRouters(r, db)
	r.Run()
}
