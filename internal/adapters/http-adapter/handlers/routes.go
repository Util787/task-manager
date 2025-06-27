package handlers

import (
	"github.com/Util787/task-manager/internal/adapters/http-adapter/handlers/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (h *Handlers) InitRoutes(env string) *gin.Engine {
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	if env != "prod" {
		router.Use(gin.Logger())
	}

	router.Use(middleware.LoggingMiddleware(h.log))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			tasks := v1.Group("/tasks")
			{
				tasks.POST("/", h.createTask)
				tasks.GET("/:id/state", h.getTaskStateByID)
				tasks.GET("/:id/result", h.getTaskResultByID)
				tasks.DELETE("/:id", h.deleteTask)
			}
		}
	}

	return router
}
