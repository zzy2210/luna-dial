package middleware

import (
	"net/http"
	"strings"

	"okr-web/internal/types"

	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware 统一的认证中间件
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 跳过不需要认证的路径
			path := c.Request().URL.Path
			if shouldSkipAuth(path) {
				return next(c)
			}

			// 获取用户ID
			userIDStr := c.Get("user_id")
			// 调试日志：输出 user_id 及类型
			log.Printf("[AuthMiddleware] user_id: %v, type: %T", userIDStr, userIDStr)
			if userIDStr == nil {
				return c.JSON(http.StatusUnauthorized, types.ErrorResponse{
					Success: false,
					Error:   "UNAUTHORIZED",
					Message: "未授权访问1",
				})
			}

			// 解析用户ID
			userID, err := uuid.Parse(userIDStr.(string))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, types.ErrorResponse{
					Success: false,
					Error:   "INVALID_USER_ID",
					Message: "用户ID格式错误",
				})
			}

			// 将解析后的UUID存入context
			c.Set("parsed_user_id", userID)
			// 调试日志：输出 parsed_user_id 及类型
			log.Printf("[AuthMiddleware] parsed_user_id: %v, type: %T", c.Get("parsed_user_id"), c.Get("parsed_user_id"))

			return next(c)
		}
	}
}

// shouldSkipAuth 判断是否跳过认证的路径
func shouldSkipAuth(path string) bool {
	// 不需要认证的路径列表
	publicPaths := []string{
		"/health",
		"/api/ping",
		"/api/auth/login",
		"/api/auth/register",
		"/api/error-test",
	}

	for _, publicPath := range publicPaths {
		if path == publicPath {
			return true
		}
	}

	// 其他所有 /api/ 开头的路径都需要认证
	return !strings.HasPrefix(path, "/api/")
}
