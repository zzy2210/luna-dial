package service

import (
	"fmt"
	"luna_dial/internal/biz"
	"time"

	"github.com/labstack/echo/v4"
)

// handleGetCurrentUser 获取当前用户详细信息
func (s *Service) handleGetCurrentUser(c echo.Context) error {
	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, fmt.Sprintf("Unauthorized: %v", err)))
	}

	session, flag := GetSessionFromContext(c)
	if !flag {
		return c.JSON(401, NewErrorResponse(401, "Unauthorized"))
	}

	// 从数据库获取用户完整信息
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
