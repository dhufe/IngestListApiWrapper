package http

import (
	"github.com/gin-gonic/gin"

	"github.com/dhufe/IngestListApiWrapper/internal/application/services"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/http/handlers"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/http/middleware"
)

func NewRouter(
	authService *services.AuthService,
	taskService *services.TaskService,
) *gin.Engine {
	router := gin.Default()

	// middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)
	// Handler initialisieren
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// public endpoints
	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
	}
	// potected endpoints

	authenticated := api.Group("")
	authenticated.Use(authMiddleware.RequireAuth())
	{
		authenticated.POST("/create", taskHandler.CreateTask)
		authenticated.GET("/", taskHandler.DefaultReponse)
		authenticated.GET("/jobs", taskHandler.GetAllTasks)
		authenticated.GET("/job/:id", taskHandler.GetTask)
		authenticated.PUT("/job/:id", taskHandler.UpdateTask)
		authenticated.DELETE("job/:id", taskHandler.DeleteTask)
	}
	return router
}
