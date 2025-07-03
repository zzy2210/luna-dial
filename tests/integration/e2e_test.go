package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"okr-web/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestMainEndpoints 测试主要的API端点
func TestMainEndpoints(t *testing.T) {
	// 创建Echo应用
	app := echo.New()

	// 添加基本的ping端点
	app.GET("/api/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	// 测试ping端点
	t.Run("Ping endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "pong", response["message"])
	})
}

// TestConfigLoad 测试配置加载
func TestConfigLoad(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Server.Port)
	assert.NotEmpty(t, cfg.Database.Host)
}

// TestRouteSetup 测试路由设置
func TestRouteSetup(t *testing.T) {
	// 创建Echo应用
	app := echo.New()

	// 模拟身份验证路由
	auth := app.Group("/api/auth")
	auth.POST("/register", func(c echo.Context) error {
		var req map[string]string
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// 简单验证
		if req["username"] == "" || req["password"] == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and password required"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User registered successfully"})
	})

	auth.POST("/login", func(c echo.Context) error {
		var req map[string]string
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// 简单验证
		if req["username"] == "" || req["password"] == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and password required"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"token":   "test-jwt-token",
			"user_id": "test-user-id",
			"message": "Login successful",
		})
	})

	// 测试注册端点
	t.Run("Register endpoint", func(t *testing.T) {
		registerData := map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "testpass123",
		}

		jsonData, err := json.Marshal(registerData)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User registered successfully", response["message"])
	})

	// 测试登录端点
	t.Run("Login endpoint", func(t *testing.T) {
		loginData := map[string]string{
			"username": "testuser",
			"password": "testpass123",
		}

		jsonData, err := json.Marshal(loginData)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user_id")
		assert.Equal(t, "Login successful", response["message"])
	})

	// 测试无效请求
	t.Run("Invalid register request", func(t *testing.T) {
		invalidData := map[string]string{
			"username": "",
			"password": "",
		}

		jsonData, err := json.Marshal(invalidData)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Username and password required", response["error"])
	})
}
