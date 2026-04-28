package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

// @Summary      Healthcheck
// @Description  Проверка работоспособности сервера
// @Tags         health
// @Success      200 {object} map[string]string "status: OK"
// @Router       /helthcheck [get]
func healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func createRouters(r *gin.Engine, db *gorm.DB) {

	r.GET("/helthcheck", func(ctx *gin.Context) { healthcheck(ctx) })

	r.GET("/services", func(c *gin.Context) { handleListSeriveces(c, db) })
	r.GET("/services/:id", func(c *gin.Context) { handleReadService(c, db) })
	r.POST("/services", func(c *gin.Context) { handleCreateService(c, db) })
	r.PUT("/services/:id", func(c *gin.Context) { handleUpdateService(c, db) })
	r.DELETE("/services/:id", func(c *gin.Context) { handleDeleteService(c, db) })
	r.GET("/services/sum", func(c *gin.Context) { handleSumService(c, db) })
}
