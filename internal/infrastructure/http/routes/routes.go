package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/http/handlers"
)

func SetupTaskRoutes(router *gin.Engine, handler *handlers.TaskHandler) {
	taskRoutes := router.Group("/api")
	{
		taskRoutes.POST("/create", handler.CreateTask)
		taskRoutes.GET("/", handler.DefaultReponse)

		taskRoutes.GET("/job/:id", handler.GetTask)
		taskRoutes.PUT("/job/:id", handler.UpdateTask)
		taskRoutes.DELETE("job/:id", handler.DeleteTask)
	}
}
