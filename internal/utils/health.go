package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"okr-web/internal/config"
	"okr-web/internal/types"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	db *config.Database
}

// NewHealthChecker 创建新的健康检查器
func NewHealthChecker(db *config.Database) *HealthChecker {
	return &HealthChecker{db: db}
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status   string                 `json:"status"`
	Service  string                 `json:"service"`
	Version  string                 `json:"version"`
	Time     string                 `json:"time"`
	Database map[string]interface{} `json:"database"`
}

// Check 执行健康检查
func (h *HealthChecker) Check(c echo.Context) error {
	response := HealthCheckResponse{
		Status:  "ok",
		Service: "OKR年度计划管理系统",
		Version: "v1.0.0",
		Time:    time.Now().Format(time.RFC3339),
		Database: map[string]interface{}{
			"status": "ok",
		},
	}

	// 检查数据库连接
	if h.db != nil {
		if err := h.db.Health(); err != nil {
			response.Status = "degraded"
			response.Database = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
			return c.JSON(http.StatusServiceUnavailable, types.ErrorResponse{
				Success: false,
				Error:   "SERVICE_UNAVAILABLE",
				Message: "数据库连接失败",
			})
		}
	} else {
		response.Database = map[string]interface{}{
			"status": "not_configured",
		}
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    response,
	})
}
