package services

import (
	"context"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/models"
)

type TaskService struct {
	repo            interfaces.TaskRepository
	fileStoragePath string
}

func NewTaskService(repo interfaces.TaskRepository, fileStoragePath string) *TaskService {
	return &TaskService{
		repo:            repo,
		fileStoragePath: fileStoragePath,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, filename string) (*models.Task, error) {
	task := &models.Task{
		FileName: filename,
		Status:   models.StatusPending,
	}

	err := s.repo.Create(ctx, task)
	return task, err
}

func (s *TaskService) FileStoragePath() string {
	return s.fileStoragePath
}

func (s *TaskService) FindPendingTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := s.repo.FindPendingTasks(ctx, &tasks)
	return tasks, err
}

func (s *TaskService) GetTask(ctx context.Context, taskID uint) (*models.Task, error) {
	return s.repo.FindByID(ctx, taskID)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID uint) error {
	return s.repo.Delete(ctx, taskID)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.repo.FindAll(ctx)
	return tasks, err
}

func (s *TaskService) UpdateTask(
	ctx context.Context,
	taskId uint,
	filename string,
	status models.TaskStatus,
) (*models.Task, error) {
	task := &models.Task{
		ID:       taskId,
		FileName: filename,
		Status:   status,
	}

	err := s.repo.Update(ctx, task)
	return task, err
}
