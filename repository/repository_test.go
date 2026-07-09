package repository

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	return db
}

func TestUserRepositoryGorm_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepositoryGorm(db)

	user := &models.User{Name: "Alice", Email: "alice@example.com"}
	err := user.HashPassword("pass123")
	require.NoError(t, err)

	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, user.Id)
}

func TestUserRepositoryGorm_FindAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepositoryGorm(db)

	users, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepositoryGorm_FindAll_WithUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepositoryGorm(db)

	user1 := &models.User{Name: "Alice", Email: "alice@example.com"}
	user1.HashPassword("pass123")
	repo.Create(context.Background(), user1)

	user2 := &models.User{Name: "Bob", Email: "bob@example.com"}
	user2.HashPassword("pass456")
	repo.Create(context.Background(), user2)

	users, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserRepositoryGorm_FindByID_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepositoryGorm(db)

	user := &models.User{Name: "Alice", Email: "alice@example.com"}
	user.HashPassword("pass123")
	repo.Create(context.Background(), user)

	found, err := repo.FindByID(context.Background(), user.Id)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", found.Name)
}

func TestUserRepositoryGorm_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepositoryGorm(db)

	_, err := repo.FindByID(context.Background(), 999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserRepositoryGorm_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepositoryGorm(db)

	user := &models.User{Name: "Alice", Email: "alice@example.com"}
	user.HashPassword("pass123")
	repo.Create(context.Background(), user)

	user.Name = "Alice Updated"
	err := repo.Update(context.Background(), user)
	assert.NoError(t, err)

	found, _ := repo.FindByID(context.Background(), user.Id)
	assert.Equal(t, "Alice Updated", found.Name)
}

func TestUserRepositoryGorm_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepositoryGorm(db)

	user := &models.User{Name: "Alice", Email: "alice@example.com"}
	user.HashPassword("pass123")
	repo.Create(context.Background(), user)

	err := repo.Delete(context.Background(), user)
	assert.NoError(t, err)

	_, err = repo.FindByID(context.Background(), user.Id)
	assert.Error(t, err)
}

func TestMockUserRepository_FindAll(t *testing.T) {
	mock := NewMockUserRepository()

	users, err := mock.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, users)

	mock.Create(context.Background(), &models.User{Name: "Alice", Email: "alice@example.com"})
	users, err = mock.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestMockUserRepository_FindByID_Found(t *testing.T) {
	mock := NewMockUserRepository()
	mock.Create(context.Background(), &models.User{Name: "Alice", Email: "alice@example.com"})

	user, err := mock.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
}

func TestMockUserRepository_FindByID_NotFound(t *testing.T) {
	mock := NewMockUserRepository()

	_, err := mock.FindByID(context.Background(), 999)
	assert.Error(t, err)
}

func TestMockUserRepository_Create(t *testing.T) {
	mock := NewMockUserRepository()

	user := &models.User{Name: "Alice", Email: "alice@example.com"}
	err := mock.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.Id)
	assert.Len(t, mock.Users, 1)
}

func TestMockUserRepository_Update(t *testing.T) {
	mock := NewMockUserRepository()
	mock.Create(context.Background(), &models.User{Name: "Alice", Email: "alice@example.com"})

	user := mock.Users[0]
	user.Name = "Updated"
	err := mock.Update(context.Background(), &user)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", mock.Users[0].Name)
}

func TestMockUserRepository_Delete(t *testing.T) {
	mock := NewMockUserRepository()
	mock.Create(context.Background(), &models.User{Name: "Alice", Email: "alice@example.com"})

	err := mock.Delete(context.Background(), &mock.Users[0])
	assert.NoError(t, err)
	assert.Len(t, mock.Users, 0)
}

func TestMockUserRepository_FindErr(t *testing.T) {
	mock := NewMockUserRepository()
	mock.FindErr = assert.AnError

	_, err := mock.FindAll(context.Background())
	assert.Error(t, err)

	_, err = mock.FindByID(context.Background(), 1)
	assert.Error(t, err)
}

func TestMockUserRepository_CreateErr(t *testing.T) {
	mock := NewMockUserRepository()
	mock.CreateErr = assert.AnError

	err := mock.Create(context.Background(), &models.User{})
	assert.Error(t, err)
}

func TestMockUserRepository_UpdateErr(t *testing.T) {
	mock := NewMockUserRepository()
	mock.Create(context.Background(), &models.User{Name: "Alice", Email: "alice@example.com"})
	mock.UpdateErr = assert.AnError

	err := mock.Update(context.Background(), &mock.Users[0])
	assert.Error(t, err)
}

func TestMockUserRepository_DeleteErr(t *testing.T) {
	mock := NewMockUserRepository()
	mock.Create(context.Background(), &models.User{Name: "Alice", Email: "alice@example.com"})
	mock.DeleteErr = assert.AnError

	err := mock.Delete(context.Background(), &mock.Users[0])
	assert.Error(t, err)
}
