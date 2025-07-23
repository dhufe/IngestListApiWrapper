package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dhufe/IngestListApiWrapper/internal/application/services"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/http/handlers"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/persistence"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/scheduler"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/worker"
	"github.com/dhufe/IngestListApiWrapper/internal/interfaces/http"
	"github.com/dhufe/IngestListApiWrapper/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config file: %v\n", err)
	}

	db, err := persistence.NewDatabase(
		cfg.Database.Driver,
		cfg.Database.DSN,
		cfg.Database.MaxOpen,
		cfg.Database.MaxIdle,
	)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v\n", err)
	}

	// Repository erstellen
	taskRepo := persistence.NewGormTaskRepository(db)
	// Service erstellen
	taskService := services.NewTaskService(taskRepo)

	// Worker erstellen (max. 3 parallele Tasks)
	taskWorker := worker.NewTaskWorker(taskRepo, 3)

	// Scheduler starten
	taskScheduler := scheduler.NewTaskScheduler(taskService, taskWorker)
	taskScheduler.Start()
	defer taskScheduler.Stop()

	// HTTP-Handler erstellen
	taskHandler := handlers.NewTaskHandler(taskService)

	// HTTP-Server erstellen und starten
	server := http.NewServer(taskHandler)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Auf Shutdown-Signal warten
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Server herunterfahren
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server exited")
}
