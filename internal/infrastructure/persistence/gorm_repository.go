package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/models"
)

type GormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) interfaces.TaskRepository {
	return &GormTaskRepository{db: db}
}

func (r *GormTaskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *GormTaskRepository) FindByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).First(&task, id).Error
	return &task, err
}

func (r *GormTaskRepository) FindAll(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).Find(&tasks).Error
	return tasks, err
}

func (r *GormTaskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *GormTaskRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, id).Error
}

func (r *GormTaskRepository) FindDueTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("status = ? AND due_date <= ?", models.StatusPending, time.Now()).
		Find(&tasks).Error
	return tasks, err
}
