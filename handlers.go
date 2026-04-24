package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ServiceCreateUpdate struct {
	ServiceName string          `json:"service_name" binding:"required"`
	Price       decimal.Decimal `json:"price" binding:"required"`
	UserID      uuid.UUID       `json:"user_id" binding:"required,uuid"`
	StartedAt   time.Time       `json:"started_at"`
	FinishedAt  time.Time       `json:"finished_at"`
}

func handleListSeriveces(ctx *gin.Context, db *gorm.DB) {
	var services []Service
	if err := db.Find(&services).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch services"})
	}
	ctx.JSON(http.StatusOK, gin.H{"services": services})
}

func handleCreateService(ctx *gin.Context, db *gorm.DB) {
	var createService ServiceCreateUpdate
	if err := ctx.ShouldBindJSON(&createService); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
	v := validator.New()
	if err := v.Struct(&createService); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	service := Service{
		ServiceName: createService.ServiceName,
		Price:       createService.Price,
		UserID:      createService.UserID,
		StartedAt:   createService.StartedAt,
		FinishedAt:  createService.FinishedAt,
	}

	if err := db.Create(&service).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": service.ID})
}

func handleReadService(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var service Service
	if err := db.First(service, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
	}
	ctx.JSON(http.StatusOK, service)
}

func handleDeleteService(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	if err := db.Delete(&Service{}, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
	}
	ctx.Status(http.StatusNoContent)
}

func handleUpdateService(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var service_update ServiceCreateUpdate
	if err := ctx.ShouldBindJSON(&service_update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
	var service Service
	if err := db.First(service, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
	}
	service.Price = service_update.Price
	service.ServiceName = service_update.ServiceName
	service.StartedAt = service_update.StartedAt
	service.FinishedAt = service_update.FinishedAt
	service.UserID = service_update.UserID
	if err := db.Save(&service).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update service"})
		return
	}
}
