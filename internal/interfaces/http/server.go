package http

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/http/handlers"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/http/routes"
)

type Server struct {
	router *gin.Engine
	server *http.Server
}

func NewServer(taskHandler *handlers.TaskHandler) *Server {
	router := gin.Default()

	// Routen registrieren
	routes.SetupTaskRoutes(router, taskHandler)

	return &Server{
		router: router,
		server: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	log.Println("Starting HTTP server on :8080")
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
