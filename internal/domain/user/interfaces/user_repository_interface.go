package repositories

import (
	"context"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/user/models"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
}
