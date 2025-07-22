package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/dhufe/IngestListApiWrapper/internal/application/services"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/models"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/http/handlers"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/persistence"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/scheduler"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/worker"
	"github.com/dhufe/IngestListApiWrapper/internal/interfaces/http"
)

func main() {
	// Datenbank initialisieren
	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Datenbank migrieren
	if err := db.AutoMigrate(&models.Task{}, &models.TaskOutput{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Repository erstellen
	taskRepo := persistence.NewGormTaskRepository(db)
	outputRepo := persistence.NewGormTaskOutputRepository(db)

	// Service erstellen
	taskService := services.NewTaskService(taskRepo, outputRepo)

	// Worker erstellen (max. 3 parallele Tasks)
	taskWorker := worker.NewTaskWorker(taskRepo, outputRepo, 3)

	// Scheduler starten
	taskScheduler := scheduler.NewTaskScheduler(taskService, taskWorker)
	taskScheduler.Start()
	defer taskScheduler.Stop()

	// HTTP-Handler erstellen
	taskHandler := handlers.NewTaskHandler(taskService, outputRepo)

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
