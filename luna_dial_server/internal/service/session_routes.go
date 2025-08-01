package service

import (
	"github.com/labstack/echo/v4"
)

// setupSessionRoutes 设置Session相关的路由
func (s *Service) setupSessionRoutes() {

	// 受保护的路由 - 需要Session认证
	protected := s.e.Group("/api/v1")
	protected.Use(s.SessionMiddleware())

	// 用户相关接口
	protected.GET("/auth/profile", s.handleGetProfile)
	protected.POST("/auth/logout", s.handleSessionLogout)
	protected.DELETE("/auth/logout-all", s.handleLogoutAllSessions)

	// 其他业务接口...
	userGroup := protected.Group("/users")
	userGroup.GET("/me", s.handleGetCurrentUser)

	journalGroup := protected.Group("/journals")
	journalGroup.GET("", s.handleListJournals)
	journalGroup.POST("", s.handleCreateJournal)

	taskGroup := protected.Group("/tasks")
	taskGroup.GET("", s.handleListTasks)
	taskGroup.POST("", s.handleCreateTask)
}

// 示例处理器方法

// handleSessionLogin 用户登录处理器
func (s *Service) handleSessionLogin(c echo.Context) error {
	// 这里应该验证用户名密码
	// 简化示例，假设验证成功
	userID := int64(1)
	username := "test_user"

	// 创建session
	sessionResp, err := s.CreateSession(c.Request().Context(), userID, username)
	if err != nil {
		return echo.NewHTTPError(500, map[string]string{
			"error":   "internal_error",
			"message": "Failed to create session",
		})
	}

	return c.JSON(200, map[string]interface{}{
		"message":    "Login successful",
		"session_id": sessionResp.SessionID,
		"expires_in": sessionResp.ExpiresIn,
	})
}

// handleSessionLogout 用户退出处理器
func (s *Service) handleSessionLogout(c echo.Context) error {
	sessionID, ok := GetSessionIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(400, map[string]string{
			"error":   "bad_request",
			"message": "Session ID not found",
		})
	}

	err := s.DeleteSession(c.Request().Context(), sessionID)
	if err != nil {
		return echo.NewHTTPError(500, map[string]string{
			"error":   "internal_error",
			"message": "Failed to logout",
		})
	}

	return c.JSON(200, map[string]string{
		"message": "Logout successful",
	})
}

// handleLogoutAllSessions 退出所有设备
func (s *Service) handleLogoutAllSessions(c echo.Context) error {
	userID, _, ok := GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(400, map[string]string{
			"error":   "bad_request",
			"message": "User ID not found",
		})
	}

	err := s.DeleteUserSessions(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(500, map[string]string{
			"error":   "internal_error",
			"message": "Failed to logout from all devices",
		})
	}

	return c.JSON(200, map[string]string{
		"message": "Logged out from all devices",
	})
}

// handleGetProfile 获取用户资料
func (s *Service) handleGetProfile(c echo.Context) error {
	userID, username, ok := GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(401, map[string]string{
			"error":   "unauthorized",
			"message": "User information not found",
		})
	}

	return c.JSON(200, map[string]interface{}{
		"user_id":  userID,
		"username": username,
	})
}

// handleGetCurrentUser 获取当前用户信息
func (s *Service) handleGetCurrentUser(c echo.Context) error {
	session, ok := GetSessionFromContext(c)
	if !ok {
		return echo.NewHTTPError(401, map[string]string{
			"error":   "unauthorized",
			"message": "Session not found",
		})
	}

	return c.JSON(200, map[string]interface{}{
		"user_id":        session.UserID,
		"username":       session.Username,
		"last_access_at": session.LastAccessAt,
		"expires_at":     session.ExpiresAt,
	})
}

// 占位符处理器（需要根据业务逻辑实现）
func (s *Service) handleListJournals(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "list journals endpoint"})
}

func (s *Service) handleCreateJournal(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "create journal endpoint"})
}

func (s *Service) handleListTasks(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "list tasks endpoint"})
}

func (s *Service) handleCreateTask(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "create task endpoint"})
}
