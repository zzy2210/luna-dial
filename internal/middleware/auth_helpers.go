package middleware

import (
	"net/http"

	"okr-web/internal/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetUserID 从上下文中获取用户ID的辅助函数
// 统一处理不同类型的用户ID并返回UUID格式
func GetUserID(c echo.Context) (uuid.UUID, error) {
	userIDValue := c.Get("user_id")
	if userIDValue == nil {
		return uuid.UUID{}, &types.AppError{
			Code:    http.StatusUnauthorized,
			Message: "未授权访问",
			Type:    "UNAUTHORIZED",
		}
	}

	// 处理不同类型的用户ID
	switch v := userIDValue.(type) {
	case string:
		userID, err := uuid.Parse(v)
		if err != nil {
			return uuid.UUID{}, &types.AppError{
				Code:    http.StatusUnauthorized,
				Message: "用户ID格式错误",
				Type:    "INVALID_USER_ID",
			}
		}
		return userID, nil
	case uuid.UUID:
		return v, nil
	default:
		return uuid.UUID{}, &types.AppError{
			Code:    http.StatusUnauthorized,
			Message: "用户ID格式错误",
			Type:    "INVALID_USER_ID",
		}
	}
}

// GetUserIDOrError 获取用户ID，如果失败则返回HTTP错误响应
func GetUserIDOrError(c echo.Context) (uuid.UUID, error) {
	userID, err := GetUserID(c)
	if err != nil {
		if appErr, ok := err.(*types.AppError); ok {
			return uuid.UUID{}, c.JSON(appErr.Code, types.ErrorResponse{
				Success: false,
				Error:   appErr.Type,
				Message: appErr.Message,
			})
		}
		return uuid.UUID{}, c.JSON(http.StatusUnauthorized, types.ErrorResponse{
			Success: false,
			Error:   "UNAUTHORIZED",
			Message: "未授权访问",
		})
	}
	return userID, nil
}
