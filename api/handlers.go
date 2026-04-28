package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TODO: add date validation
type ServiceCreateUpdate struct {
	ServiceName string          `json:"service_name" binding:"required"`
	Price       decimal.Decimal `json:"price" binding:"required"`
	UserID      uuid.UUID       `json:"user_id" binding:"required,uuid"`
	StartedAt   string          `json:"started_at" binding:"required" time_format:"01-2026"`
	FinishedAt  string          `json:"finished_at" time_format:"01-2006"`
}

type ServiceRead struct {
	ID          uuid.UUID       `json:"id" binding:"required"`
	ServiceName string          `json:"service_name" binding:"required"`
	Price       decimal.Decimal `json:"price" binding:"required"`
	UserID      uuid.UUID       `json:"user_id" binding:"required,uuid"`
	StartedAt   time.Time       `json:"started_at" binding:"required" time_format:"2026-01-01"`
	FinishedAt  time.Time       `json:"finished_at" binding:"required" time_format:"2026-01-01"`
}

type SumResult struct {
	ServiceName string          `json:"service_name" binding:"required"`
	Price       decimal.Decimal `json:"price" binding:"required"`
	UserID      uuid.UUID       `json:"user_id" binding:"uuid"`
}

func handleListSeriveces(ctx *gin.Context, db *gorm.DB) {
	var services []Service
	var results []ServiceRead
	var offset int
	var limit int
	var count int64
	if ctx.Query("offset") == "" {
		offset = 0
	} else {
		_offset, err := strconv.Atoi(ctx.Query("offset"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid offset value"})
			return
		}
		offset = _offset
	}
	if ctx.Query("limit") == "" {
		limit = 10
	} else {
		_limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid limit value"})
			return
		}
		limit = _limit
	}
	if err := db.Limit(limit).Offset(offset).Find(&services).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch services"})
		return
	}
	for _, s := range services {
		results = append(results, ServiceRead{
			ID:          s.ID,
			ServiceName: s.ServiceName,
			Price:       s.Price,
			UserID:      s.UserID,
			StartedAt:   s.StartedAt,
			FinishedAt:  s.FinishedAt,
		})
	}
	db.Model(&services).Count(&count)

	ctx.JSON(http.StatusOK, gin.H{"data": results, "count": count, "limit": limit, "offset": offset})
}

func handleCreateService(ctx *gin.Context, db *gorm.DB) {
	var createService ServiceCreateUpdate
	if err := ctx.ShouldBindJSON(&createService); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	v := validator.New()
	if err := v.Struct(&createService); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	st, err := time.Parse("01-2006", createService.StartedAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid started_at"})
		return
	}
	print(st.String())
	ft, err := time.Parse("01-2006", createService.FinishedAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid finished_at"})
		return
	}

	service := Service{
		ServiceName: createService.ServiceName,
		Price:       createService.Price,
		UserID:      createService.UserID,
		StartedAt:   time.Date(st.Year(), st.Month(), 1, 0, 0, 0, 0, time.UTC),
		FinishedAt:  time.Date(ft.Year(), ft.Month()+1, 0, 23, 59, 59, 0, time.UTC),
	}

	if err := db.Create(&service).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": service.ID})
}

func handleReadService(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var service Service
	if err := db.Where("id = ?", id).First(&service).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	ctx.JSON(http.StatusOK, &ServiceRead{
		ID:          service.ID,
		ServiceName: service.ServiceName,
		Price:       service.Price,
		UserID:      service.UserID,
		StartedAt:   service.StartedAt,
		FinishedAt:  service.FinishedAt,
	})
}

func handleDeleteService(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var service Service
	if err := db.Where("id = ?", id).Delete(&service).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func handleUpdateService(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var service_update ServiceCreateUpdate
	if err := ctx.ShouldBindJSON(&service_update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	var service Service
	if err := db.First(&service, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	st, err := time.Parse("01-2006", service_update.StartedAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid started_at"})
		return
	}
	ft, err := time.Parse("01-2006", service_update.FinishedAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid finished_at"})
		return
	}
	service.Price = service_update.Price
	service.ServiceName = service_update.ServiceName
	service.StartedAt = time.Date(st.Year(), st.Month(), 1, 0, 0, 0, 0, time.UTC)
	service.FinishedAt = time.Date(ft.Year(), ft.Month()+1, 0, 23, 59, 59, 99, time.UTC)
	service.UserID = service_update.UserID
	if err := db.Save(&service).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update service"})
		return
	}
}

func handleSumService(ctx *gin.Context, db *gorm.DB) {
	var result SumResult
	user_id := ctx.Query("user_id")
	if user_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "need user_id"})
		return
	}
	service_name := ctx.Query("service_name")
	if service_name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "need service_name"})
		return
	}
	start, err := time.Parse("01-2006", ctx.Query("start"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start"})
		return
	}
	finish, err := time.Parse("01-2006", ctx.Query("finish"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid finish"})
		return
	}

	db.Model(&Service{}).
		Select("service_name, sum(price) as price, user_id").
		Where("user_id = ?", user_id).
		Where("service_name = ?", service_name).
		Where("started_at BETWEEN ? AND ?", start, finish).
		Group("service_name, user_id").
		Order("service_name").
		Find(&result)
	ctx.JSON(http.StatusOK, result)
}
