package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/models"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) interfaces.TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *TaskRepository) FindByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).First(&task, id).Error
	return &task, err
}

func (r *TaskRepository) FindAll(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *TaskRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, id).Error
}

func (r *TaskRepository) FindDueTasks(ctx context.Context, tasks *[]models.Task) error {
	err := r.db.WithContext(ctx).
		Where("status = ?", models.StatusPending).
		Find(tasks).Error
	return err
}
