package services

import (
	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/interfaces"
)

type MetricsService struct {
	repo            interfaces.TaskRepository
	fileStoragePath string
}

func NewMetricsService(repo interfaces.TaskRepository, fileStoragePath string) *MetricsService {
	return &MetricsService{
		repo:            repo,
		fileStoragePath: fileStoragePath,
	}
}
