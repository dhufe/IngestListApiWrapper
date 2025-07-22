package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/models"
)

type GormTaskOutputRepository struct {
	db *gorm.DB
}

func NewGormTaskOutputRepository(db *gorm.DB) interfaces.TaskOutputRepository {
	return &GormTaskOutputRepository{db: db}
}

func (r *GormTaskOutputRepository) Create(ctx context.Context, output *models.TaskOutput) error {
	return r.db.WithContext(ctx).Create(output).Error
}

func (r *GormTaskOutputRepository) FindByTaskID(
	ctx context.Context,
	taskID uint,
) ([]models.TaskOutput, error) {
	var outputs []models.TaskOutput
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&outputs).Error
	return outputs, err
}

func (r *GormTaskOutputRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TaskOutput{}, id).Error
}

func (r *GormTaskOutputRepository) DeleteByTaskID(ctx context.Context, taskID uint) error {
	return r.db.WithContext(ctx).Where("task_id = ?", taskID).Delete(&models.TaskOutput{}).Error
}
