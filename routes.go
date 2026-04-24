package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func createRouters(r *gin.Engine, db *gorm.DB) {
	r.GET("/helthcheck", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "OK"}) })

	r.GET("/services", func(c *gin.Context) { handleListSeriveces(c, db) })
	r.GET("/services/:id", func(c *gin.Context) { handleReadService(c, db) })
	r.POST("/services", func(c *gin.Context) { handleCreateService(c, db) })
	r.PUT("/services/:id", func(c *gin.Context) { handleUpdateService(c, db) })
	r.DELETE("/services/:id", func(c *gin.Context) { handleListSeriveces(c, db) })
}
