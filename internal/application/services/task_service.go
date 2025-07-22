package services

import (
	"context"
	"time"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/models"
)

type TaskService struct {
	repo interfaces.TaskRepository
}

func NewTaskService(repo interfaces.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	title, description string,
	dueDate *time.Time,
) (*models.Task, error) {
	task := &models.Task{
		Title: title,

		Status:  models.StatusPending,
		DueDate: dueDate,
	}

	err := s.repo.Create(ctx, task)
	return task, err
}

func (s *TaskService) GetTask(ctx context.Context, id uint) (*models.Task, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.repo.FindAll(ctx)
}

func (s *TaskService) UpdateTask(
	ctx context.Context,
	id uint,
	title, description string,
	status models.TaskStatus,
	dueDate *time.Time,
) (*models.Task, error) {
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Title = title

	task.Status = status
	task.DueDate = dueDate

	err = s.repo.Update(ctx, task)
	return task, err
}

func (s *TaskService) DeleteTask(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *TaskService) ProcessDueTasks(ctx context.Context) error {
	tasks := []models.Task{}
	err := s.repo.FindDueTasks(ctx, &tasks)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		task.Status = models.StatusPending
		if err := s.repo.Update(ctx, &task); err != nil {
			return err
		}
	}

	return nil
}
