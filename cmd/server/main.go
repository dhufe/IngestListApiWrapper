package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dhufe/IngestListApiWrapper/internal/application/services"

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
	taskRepo := persistence.NewTaskRepository(db)
	userRepo := persistence.NewUserRepository(db)
	// Service erstellen
	taskService := services.NewTaskService(taskRepo)

	// Worker erstellen (max. Anzahl v. parallelen Tasks)
	taskWorker := worker.NewTaskWorker(taskRepo, cfg.TaskScheduler.MaxWorkers)

	// Scheduler starten
	taskScheduler := scheduler.NewTaskScheduler(taskService, taskWorker, cfg.TaskScheduler.Interval)
	taskScheduler.Start()
	defer taskScheduler.Stop()

	authService := services.NewAuthService(
		userRepo,
		"your-secret-key-for-jwt",
		2*time.Hour,
	)

	// HTTP-Server erstellen und starten
	// server := http.NewServer(taskHandler)
	server := http.NewRouter(authService, taskService)
	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Auf Shutdown-Signal warten
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	log.Println("Server exited")
}
