package services

import (
	"context"
	"time"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/models"
)

type TaskService struct {
	repo       interfaces.TaskRepository
	outputRepo interfaces.TaskOutputRepository
}

func NewTaskService(repo interfaces.TaskRepository, outputRepo interfaces.TaskOutputRepository) *TaskService {
	return &TaskService{
		repo:       repo,
		outputRepo: outputRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title, command, arguments string, dueDate *time.Time) (*models.Task, error) {
	task := &models.Task{
		Title:     title,
		Command:   command,
		Arguments: arguments,
		Status:    models.StatusPending,
		DueDate:   dueDate,
	}

	err := s.repo.Create(ctx, task)
	return task, err
}

func (s *TaskService) FindDueTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := s.repo.FindDueTasks(ctx, &tasks)
	return tasks, err
}

func (s *TaskService) GetTaskOutput(ctx context.Context, taskID uint) ([]models.TaskOutput, error) {
	return s.outputRepo.FindByTaskID(ctx, taskID)
}

func (s *TaskService) DeleteTaskOutput(ctx context.Context, outputID uint) error {
	return s.outputRepo.Delete(ctx, outputID)
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

func (s *TaskService) UpdateTask(ctx context.Context, taskId uint, title, command, arguments string, status models.TaskStatus, dueDate *time.Time) (*models.Task, error) {
	task := &models.Task{
		ID:        taskId,
		Title:     title,
		Command:   command,
		Arguments: arguments,
		Status:    status,
		DueDate:   dueDate,
	}

	err := s.repo.Update(ctx, task)
	return task, err
}
