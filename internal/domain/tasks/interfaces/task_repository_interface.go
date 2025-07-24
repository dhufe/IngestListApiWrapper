package interfaces

import (
	"context"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	FindByID(ctx context.Context, id uint) (*models.Task, error)
	FindAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uint) error
	FindPendingTasks(ctx context.Context, tasks *[]models.Task) error
}
