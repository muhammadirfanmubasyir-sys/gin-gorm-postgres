package repository

import (
	"context"

	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	Users      []models.User
	NextID     int
	CreateErr  error
	UpdateErr  error
	DeleteErr  error
	FindErr    error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:  []models.User{},
		NextID: 1,
	}
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	if m.FindErr != nil {
		return nil, m.FindErr
	}
	result := make([]models.User, len(m.Users))
	copy(result, m.Users)
	return result, nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int) (models.User, error) {
	if m.FindErr != nil {
		return models.User{}, m.FindErr
	}
	for _, u := range m.Users {
		if u.Id == id {
			return u, nil
		}
	}
	return models.User{}, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	user.Id = m.NextID
	m.NextID++
	m.Users = append(m.Users, *user)
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	for i, u := range m.Users {
		if u.Id == user.Id {
			m.Users[i] = *user
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockUserRepository) Delete(ctx context.Context, user *models.User) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	for i, u := range m.Users {
		if u.Id == user.Id {
			m.Users = append(m.Users[:i], m.Users[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}
