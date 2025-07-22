package interfaces

import (
	"context"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/models"
)

type TaskOutputRepository interface {
	Create(ctx context.Context, output *models.TaskOutput) error
	FindByTaskID(ctx context.Context, taskID uint) ([]models.TaskOutput, error)
	Delete(ctx context.Context, id uint) error
	DeleteByTaskID(ctx context.Context, taskID uint) error
}
