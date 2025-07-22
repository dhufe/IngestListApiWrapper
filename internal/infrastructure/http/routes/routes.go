package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/http/handlers"
)

func SetupTaskRoutes(router *gin.Engine, handler *handlers.TaskHandler) {
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.POST("/", handler.CreateTask)
		taskRoutes.GET("/", handler.GetAllTasks)
		taskRoutes.GET("/:id", handler.GetTask)
		taskRoutes.PUT("/:id", handler.UpdateTask)
		taskRoutes.DELETE("/:id", handler.DeleteTask)
		taskRoutes.GET("/:id/output", handler.GetTaskOutput)
		taskRoutes.DELETE("/:id/output/:output_id", handler.DeleteTaskOutput)
	}
}

