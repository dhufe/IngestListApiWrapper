package scheduler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"github.com/dhufe/IngestListApiWrapper/internal/application/services"
	"github.com/dhufe/IngestListApiWrapper/internal/infrastructure/worker"
)

type TaskScheduler struct {
	service *services.TaskService
	worker  *worker.TaskWorker
	cron    *cron.Cron
}

func NewTaskScheduler(service *services.TaskService, worker *worker.TaskWorker) *TaskScheduler {
	return &TaskScheduler{
		service: service,
		worker:  worker,
		cron:    cron.New(),
	}
}

func (s *TaskScheduler) Start() {
	// Alle 30 Sekunden nach Tasks suchen
	_, err := s.cron.AddFunc("*/30 * * * * *", func() {
		ctx := context.Background()

		// Überfällige Tasks finden
		tasks, err := s.service.FindDueTasks(ctx)
		if err != nil {
			log.Printf("Error finding due tasks: %v", err)
			return
		}

		// Tasks verarbeiten
		for _, task := range tasks {
			s.worker.StartTaskProcessing(ctx, task)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	s.cron.Start()
}

func (s *TaskScheduler) Stop() {
	s.cron.Stop()
	s.worker.Wait()
}

