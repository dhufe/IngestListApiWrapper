package http

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"hufschlager.net/IngestListApiWrapper/internal/application/services"
	"hufschlager.net/IngestListApiWrapper/internal/infrastructure/http/handlers"
	"hufschlager.net/IngestListApiWrapper/internal/infrastructure/http/middleware"
)

func NewRouter(
	authService *services.AuthService,
	taskService *services.TaskService,
	metricsService *services.MetricsService,
) *gin.Engine {
	router := gin.Default()

	router.ForwardedByClientIP = true
	if err := router.SetTrustedProxies([]string{"*"}); err != nil {
		log.Printf("Can not set trusted proxies.\n")
	}
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))

	// middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)
	// Handler initialisieren
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)
	metricsHanlder := handlers.NewMetricsHandler(metricsService)
	/// Prometheus
	router.GET("metrics", metricsHanlder.GetMetrics)

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
