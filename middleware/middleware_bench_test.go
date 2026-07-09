package middleware

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkRequestID(b *testing.B) {
	router := setupRouter(RequestID())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkCORS(b *testing.B) {
	router := setupRouter(CORS())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://example.com")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkCORS_Preflight(b *testing.B) {
	router := setupRouter(CORS())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkSecurityHeaders(b *testing.B) {
	router := setupRouter(SecurityHeaders())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkRateLimit(b *testing.B) {
	router := setupRouter(RateLimit(1000000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkLogger(b *testing.B) {
	router := setupRouter(Logger())

	// Suppress log output during benchmark
	log.SetOutput(io.Discard)
	defer log.SetOutput(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkRecovery(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Recovery())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkErrorHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

func BenchmarkFullMiddlewareChain(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(
		RequestID(),
		CORS(),
		SecurityHeaders(),
		Logger(),
		RateLimit(1000000),
		Recovery(),
		ErrorHandler(),
	)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}
