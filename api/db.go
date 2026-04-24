package main

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func init_db(cfg *Config, log *zap.Logger) (*gorm.DB, error) {
	logger := zapgorm2.New(log)
	logger.SetAsDefault()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{Logger: logger.LogMode(gormlog.Info)})
	if err != nil {
		log.Error("Failed connect to database", zap.String("error", err.Error()))
	}
	return db, err
}

func auto_migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Service{},
	)
}
