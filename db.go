package main

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func init_db(cfg *Config, log *zap.Logger) (*gorm.DB, error) {
	logger := zapgorm2.New(log)
	logger.SetAsDefault()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{Logger: logger})
	if err != nil {
		log.Error("Can't connect to database: %v", zap.String("string", err.Error()))
	}
	return db, err
}

func auto_migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Service{},
	)
}
