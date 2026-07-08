package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"gorm.io/gorm"
)

func setupBenchDB(b *testing.B) *gorm.DB {
	b.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatalf("failed to open bench db: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		b.Fatalf("failed to migrate bench db: %v", err)
	}
	return db
}

func seedRepo(b *testing.B, repo *UserRepositoryGorm, n int) {
	b.Helper()
	for i := 0; i < n; i++ {
		user := &models.User{
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
		}
		if err := repo.Create(context.Background(), user); err != nil {
			b.Fatalf("failed to seed user: %v", err)
		}
	}
}

func BenchmarkRepo_Create(b *testing.B) {
	db := setupBenchDB(b)
	repo := NewUserRepositoryGorm(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &models.User{
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
		}
		repo.Create(context.Background(), user)
	}
}

func BenchmarkRepo_FindAll_Empty(b *testing.B) {
	db := setupBenchDB(b)
	repo := NewUserRepositoryGorm(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindAll(context.Background())
	}
}

func BenchmarkRepo_FindAll_10(b *testing.B) {
	db := setupBenchDB(b)
	repo := NewUserRepositoryGorm(db)
	seedRepo(b, repo, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindAll(context.Background())
	}
}

func BenchmarkRepo_FindByID(b *testing.B) {
	db := setupBenchDB(b)
	repo := NewUserRepositoryGorm(db)
	seedRepo(b, repo, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindByID(context.Background(), 5)
	}
}

func BenchmarkRepo_Update(b *testing.B) {
	db := setupBenchDB(b)
	repo := NewUserRepositoryGorm(db)
	seedRepo(b, repo, 1)

	user, _ := repo.FindByID(context.Background(), 1)
	ptr := &user

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ptr.Name = fmt.Sprintf("Updated%d", i)
		repo.Update(context.Background(), ptr)
	}
}

func BenchmarkRepo_Delete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		db := setupBenchDB(b)
		repo := NewUserRepositoryGorm(db)
		user := &models.User{Name: "Alice", Email: "alice@example.com"}
		repo.Create(context.Background(), user)
		b.StartTimer()

		repo.Delete(context.Background(), user)
	}
}

func BenchmarkRepo_CreateParallel(b *testing.B) {
	db := setupBenchDB(b)
	repo := NewUserRepositoryGorm(db)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			user := &models.User{
				Name:  fmt.Sprintf("User%d", i),
				Email: fmt.Sprintf("user%d@example.com", i),
			}
			repo.Create(context.Background(), user)
			i++
		}
	})
}

func BenchmarkRepo_FindAllParallel(b *testing.B) {
	db := setupBenchDB(b)
	repo := NewUserRepositoryGorm(db)
	seedRepo(b, repo, 10)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			repo.FindAll(context.Background())
		}
	})
}
