package service

import (
	"luna_dial/internal/biz"
	"time"

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
	// 解析登录请求
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, &Response{
			Code:      400,
			Message:   "Invalid request data",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	// 验证用户名密码
	u, err := s.userUsecase.UserLogin(c.Request().Context(), biz.UserLoginParam{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(401, &Response{
			Code:      401,
			Message:   "Invalid username or password",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	// 创建session
	sessionResp, err := s.CreateSession(c.Request().Context(), u.ID, u.Username)
	if err != nil {
		return c.JSON(500, &Response{
			Code:      500,
			Message:   "Failed to create session",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	return c.JSON(200, &Response{
		Code:    200,
		Message: "Login successful",
		Data: &LoginResponse{
			SessionID: sessionResp.SessionID,
			ExpiresIn: sessionResp.ExpiresIn,
		},
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

// handleSessionLogout 用户退出处理器
func (s *Service) handleSessionLogout(c echo.Context) error {
	sessionID, ok := GetSessionIDFromContext(c)
	if !ok {
		return c.JSON(400, &Response{
			Code:      400,
			Message:   "Session ID not found",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	err := s.DeleteSession(c.Request().Context(), sessionID)
	if err != nil {
		return c.JSON(500, &Response{
			Code:      500,
			Message:   "Failed to logout",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	return c.JSON(200, &Response{
		Code:      200,
		Message:   "Logout successful",
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

// handleLogoutAllSessions 退出所有设备
func (s *Service) handleLogoutAllSessions(c echo.Context) error {
	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(400, &Response{
			Code:      400,
			Message:   "User ID not found",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	err = s.DeleteUserSessions(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(500, &Response{
			Code:      500,
			Message:   "Failed to logout from all devices",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	return c.JSON(200, &Response{
		Code:      200,
		Message:   "Logged out from all devices",
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

// handleGetProfile 获取用户基本认证信息
func (s *Service) handleGetProfile(c echo.Context) error {
	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, &Response{
			Code:      401,
			Message:   "User information not found",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	// 从数据库获取用户基本信息
	user, err := s.userUsecase.GetUser(c.Request().Context(), biz.GetUserParam{
		UserID: userID,
	})
	if err != nil {
		return c.JSON(500, &Response{
			Code:      500,
			Message:   "Failed to get user information",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	return c.JSON(200, &Response{
		Code:    200,
		Message: "success",
		Data: map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
		},
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

// handleGetCurrentUser 获取当前用户详细信息
func (s *Service) handleGetCurrentUser(c echo.Context) error {
	session, ok := GetSessionFromContext(c)
	if !ok {
		return c.JSON(401, &Response{
			Code:      401,
			Message:   "Session not found",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	// 从数据库获取用户完整信息
	user, err := s.userUsecase.GetUser(c.Request().Context(), biz.GetUserParam{
		UserID: session.UserID,
	})
	if err != nil {
		return c.JSON(500, &Response{
			Code:      500,
			Message:   "Failed to get user information",
			Success:   false,
			Timestamp: time.Now().Unix(),
		})
	}

	return c.JSON(200, &Response{
		Code:    200,
		Message: "success",
		Data: map[string]interface{}{
			// 基本信息
			"user_id":  user.ID,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			
			// 账户信息
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
			
			// 会话信息
			"session": map[string]interface{}{
				"session_id":     session.ID,
				"last_access_at": session.LastAccessAt,
				"expires_at":     session.ExpiresAt,
			},
		},
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

// 占位符处理器（需要根据业务逻辑实现）
func (s *Service) handleListJournals(c echo.Context) error {
	return c.JSON(200, &Response{
		Code:      200,
		Message:   "list journals endpoint",
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

func (s *Service) handleCreateJournal(c echo.Context) error {
	return c.JSON(200, &Response{
		Code:      200,
		Message:   "create journal endpoint",
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

func (s *Service) handleListTasks(c echo.Context) error {
	return c.JSON(200, &Response{
		Code:      200,
		Message:   "list tasks endpoint",
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

func (s *Service) handleCreateTask(c echo.Context) error {
	return c.JSON(200, &Response{
		Code:      200,
		Message:   "create task endpoint",
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}
