package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ServiceName string          `gorm:"size:255"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2)"`
	UserID      uuid.UUID       `gorm:"type:uuid"`
	StartedAt   time.Time
	FinishedAt  time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
