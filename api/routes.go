package main

import (
	"net/http"

	_ "https://github.com/ruha02/go-test/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func createRouters(r *gin.Engine, db *gorm.DB) {
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// @Summary      Healthcheck
	// @Description  Проверка работоспособности сервера
	// @Tags         health
	// @Success      200 {object} map[string]string "status: OK"
	// @Router       /helthcheck [get]
	r.GET("/helthcheck", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "OK"}) })
	// @Summary      Список сервисов
	// @Description  Получить постраничный список сервисов
	// @Tags         services
	// @Accept       json
	// @Produce      json
	// @Param        offset query int false "Смещение"
	// @Param        limit  query int false "Количество"
	// @Success      200 {object} map[string]interface{} "data, count, limit, offset"
	// @Failure      500 {object} map[string]string "error"
	// @Router       /services [get]
	r.GET("/services", func(c *gin.Context) { handleListSeriveces(c, db) })
	// @Summary      Получить сервис по ID
	// @Description  Возвращает детали сервиса по его ID
	// @Tags         services
	// @Accept       json
	// @Produce      json
	// @Param        id path string true "ID сервиса"
	// @Success      200 {object} ServiceRead "service"
	// @Failure      404 {object} map[string]string "error"
	// @Router       /services/{id} [get]
	r.GET("/services/:id", func(c *gin.Context) { handleReadService(c, db) })
	// @Summary      Создать сервис
	// @Description  Создаёт новый сервис
	// @Tags         services
	// @Accept       json
	// @Produce      json
	// @Param        service body ServiceCreateUpdate true "Данные сервиса"
	// @Success      201 {object} map[string]string "id"
	// @Failure      400 {object} map[string]string "error"
	// @Router       /services [post]
	r.POST("/services", func(c *gin.Context) { handleCreateService(c, db) })
	// @Summary      Обновить сервис
	// @Description  Обновляет существующий сервис по ID
	// @Tags         services
	// @Accept       json
	// @Produce      json
	// @Param        id      path string true "ID сервиса"
	// @Param        service body ServiceCreateUpdate true "Новые данные сервиса"
	// @Success      200 {object} map[string]string "updated"
	// @Failure      400 {object} map[string]string "error"
	// @Router       /services/{id} [put]
	r.PUT("/services/:id", func(c *gin.Context) { handleUpdateService(c, db) })
	// @Summary      Удалить сервис
	// @Description  Удаляет сервис по его ID
	// @Tags         services
	// @Accept       json
	// @Produce      json
	// @Param        id path string true "ID сервиса"
	// @Success      204 "No Content"
	// @Failure      404 {object} map[string]string "error"
	// @Router       /services/{id} [delete]
	r.DELETE("/services/:id", func(c *gin.Context) { handleDeleteService(c, db) })
	// @Summary      Сумма сервисов
	// @Description  Вычисляет суммарную стоимость сервисов по пользователю и названию за указанный период
	// @Tags         services
	// @Accept       json
	// @Produce      json
	// @Param        user_id       query string true "ID пользователя"
	// @Param        service_name  query string true "Название сервиса"
	// @Param        start         query string true "Месяц начала (MM-YYYY)"
	// @Param        finish        query string true "Месяц окончания (MM-YYYY)"
	// @Success      200 {object} map[string]decimal.Decimal "sum"
	// @Failure      400 {object} map[string]string "error"
	// @Router       /services/sum [get]
	r.GET("/services/sum", func(c *gin.Context) { handleSumService(c, db) })
}
