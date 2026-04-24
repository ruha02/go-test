package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type Service struct {
	ServiceName string          `gorm:"size:255"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2)"`
	UserID      uuid.UUID       `gorm:"type:uuid,default:uuid_generate_v4(),primaryKey"`
	StartedAt   time.Time
	FinishedAt  time.Time
}

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
