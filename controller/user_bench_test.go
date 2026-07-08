package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/repository"
)

func seedBenchmarkUsers(mock *repository.MockUserRepository, n int) {
	for i := 0; i < n; i++ {
		mock.Create(nil, &models.User{
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
		})
	}
}

func BenchmarkCreateUser(b *testing.B) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)
	body := `{"name":"BenchUser","email":"bench@example.com","password":"benchpass123"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkListUser_Empty(b *testing.B) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkListUser_100(b *testing.B) {
	mock := repository.NewMockUserRepository()
	seedBenchmarkUsers(mock, 100)
	router := setupMockRouter(mock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkGetUser(b *testing.B) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	router := setupMockRouter(mock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkUpdateUser(b *testing.B) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	router := setupMockRouter(mock)
	body := `{"name":"Updated","email":"updated@example.com","password":"newpass123"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkDeleteUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mock := repository.NewMockUserRepository()
		mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
		router := setupMockRouter(mock)
		b.StartTimer()

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkCRUD_FullCycle(b *testing.B) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := fmt.Sprintf(`{"name":"User%d","email":"user%d@example.com","password":"pass123"}`, i, i)

		createReq := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(body))
		createReq.Header.Set("Content-Type", "application/json")
		createRec := httptest.NewRecorder()
		router.ServeHTTP(createRec, createReq)

		listReq := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		listRec := httptest.NewRecorder()
		router.ServeHTTP(listRec, listReq)

		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
		getRec := httptest.NewRecorder()
		router.ServeHTTP(getRec, getReq)

		updateBody := fmt.Sprintf(`{"name":"Updated%d","email":"updated%d@example.com","password":"newpass123"}`, i, i)
		updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", strings.NewReader(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
		updateRec := httptest.NewRecorder()
		router.ServeHTTP(updateRec, updateReq)

		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		deleteRec := httptest.NewRecorder()
		router.ServeHTTP(deleteRec, deleteReq)
	}
}
