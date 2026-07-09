package repository

import (
	"context"

	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]models.User, error)
	FindByID(ctx context.Context, id int) (models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
}
