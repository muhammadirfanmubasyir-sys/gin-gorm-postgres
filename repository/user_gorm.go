package repository

import (
	"context"

	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"gorm.io/gorm"
)

type UserRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepositoryGorm(db *gorm.DB) *UserRepositoryGorm {
	return &UserRepositoryGorm{db: db}
}

func (r *UserRepositoryGorm) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *UserRepositoryGorm) FindByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *UserRepositoryGorm) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryGorm) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepositoryGorm) Delete(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}
