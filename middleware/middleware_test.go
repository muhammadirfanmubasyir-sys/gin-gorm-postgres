package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(mw ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	for _, m := range mw {
		router.Use(m)
	}
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	return router
}

func performRequest(router http.Handler, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestCORS_AllowAll(t *testing.T) {
	router := setupRouter(CORS())

	rec := performRequest(router, http.MethodGet, "/test", map[string]string{
		"Origin": "http://example.com",
	})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_Preflight(t *testing.T) {
	router := setupRouter(CORS())

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCORS_SpecificOrigin(t *testing.T) {
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://allowed.com,http://also-allowed.com")
	defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

	router := setupRouter(CORS())

	rec := performRequest(router, http.MethodGet, "/test", map[string]string{
		"Origin": "http://allowed.com",
	})
	assert.Equal(t, "http://allowed.com", rec.Header().Get("Access-Control-Allow-Origin"))

	rec2 := performRequest(router, http.MethodGet, "/test", map[string]string{
		"Origin": "http://not-allowed.com",
	})
	assert.Equal(t, "", rec2.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_Headers(t *testing.T) {
	router := setupRouter(CORS())

	rec := performRequest(router, http.MethodGet, "/test", nil)

	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "X-Request-ID")
	assert.Equal(t, "86400", rec.Header().Get("Access-Control-Max-Age"))
}

func TestSecurityHeaders(t *testing.T) {
	router := setupRouter(SecurityHeaders())

	rec := performRequest(router, http.MethodGet, "/test", nil)

	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", rec.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", rec.Header().Get("Referrer-Policy"))
	assert.Equal(t, "no-store, no-cache, must-revalidate", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", rec.Header().Get("Pragma"))
}

func TestRequestID_AutoGenerate(t *testing.T) {
	router := setupRouter(RequestID())

	rec := performRequest(router, http.MethodGet, "/test", nil)

	assert.NotEmpty(t, rec.Header().Get("X-Request-ID"))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequestID_Passthrough(t *testing.T) {
	router := setupRouter(RequestID())

	customID := "my-custom-request-id"
	rec := performRequest(router, http.MethodGet, "/test", map[string]string{
		"X-Request-ID": customID,
	})

	assert.Equal(t, customID, rec.Header().Get("X-Request-ID"))
}

func TestRequestID_ContextValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())

	var contextID string
	router.GET("/test", func(c *gin.Context) {
		if val, exists := c.Get(RequestIDKey); exists {
			contextID = val.(string)
		}
		c.JSON(200, gin.H{"ok": true})
	})

	rec := performRequest(router, http.MethodGet, "/test", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, contextID)
	assert.Equal(t, rec.Header().Get("X-Request-ID"), contextID)
}

func TestRateLimit_AllowWithinLimit(t *testing.T) {
	router := setupRouter(RateLimit(5))

	for i := 0; i < 5; i++ {
		rec := performRequest(router, http.MethodGet, "/test", nil)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestRateLimit_ExceedLimit(t *testing.T) {
	router := setupRouter(RateLimit(3))

	for i := 0; i < 3; i++ {
		performRequest(router, http.MethodGet, "/test", nil)
	}

	rec := performRequest(router, http.MethodGet, "/test", nil)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "rate limit exceeded", body["error"])
	assert.Equal(t, "RATE_LIMIT_EXCEEDED", body["code"])
}

func TestRateLimit_DifferentIPs(t *testing.T) {
	router := setupRouter(RateLimit(1))

	rec1 := performRequest(router, http.MethodGet, "/test", nil)
	assert.Equal(t, http.StatusOK, rec1.Code)

	gin.SetMode(gin.TestMode)
	router2 := gin.New()
	router2.Use(RateLimit(1))
	router2.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	rec2 := httptest.NewRecorder()
	router2.ServeHTTP(rec2, req)
	assert.Equal(t, http.StatusOK, rec2.Code)
}

func TestLogger(t *testing.T) {
	router := setupRouter(Logger())

	rec := performRequest(router, http.MethodGet, "/test", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Recovery())
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	rec := performRequest(router, http.MethodGet, "/panic", nil)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestErrorHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/ok", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	rec := performRequest(router, http.MethodGet, "/ok", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestErrorHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/bad", func(c *gin.Context) {
		c.JSON(400, gin.H{"error": "bad request"})
	})

	rec := performRequest(router, http.MethodGet, "/bad", nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
