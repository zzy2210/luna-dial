package service

import (
	"luna_dial/internal/biz"
	"time"

	"github.com/labstack/echo/v4"
)

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
