package middleware

import (
	"log"
	"net/http"
	"okr-web/internal/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetUserIDFromContext 从上下文中获取用户ID（先取 string，再转 uuid.UUID）
// 这个函数统一了用户身份验证逻辑，避免在各个控制器中重复代码
// AuthMiddleware 已经处理了认证和UUID解析，直接获取parsed_user_id即可
func GetUserIDFromContext(ctx echo.Context) (uuid.UUID, error) {
	// 首先尝试获取已解析的用户ID字符串（来自AuthMiddleware）
	userID := ctx.Get("parsed_user_id")
	// 调试日志：打印 parsed_user_id 的值和类型
	log.Printf("[GetUserIDFromContext] parsed_user_id: %v, type: %T", userID, userID)
	if userID != nil {
		if uid, ok := userID.(uuid.UUID); ok {
			log.Printf("[GetUserIDFromContext] parsed_user_id 已为 uuid.UUID: %v", uid)
			return uid, nil
		}
		// 兼容老数据
		if uidStr, ok := userID.(string); ok {
			log.Printf("[GetUserIDFromContext] parsed_user_id 为 string: %v", uidStr)
			uuidVal, err := uuid.Parse(uidStr)
			if err == nil {
				log.Printf("[GetUserIDFromContext] string 成功解析为 uuid.UUID: %v", uuidVal)
				return uuidVal, nil
			} else {
				log.Printf("[GetUserIDFromContext] string 解析 uuid 失败: %v", err)
			}
		} else {
			log.Printf("[GetUserIDFromContext] parsed_user_id 类型未知: %T", userID)
		}
	}

	// 如果没有找到parsed_user_id，返回错误
	// 这通常意味着AuthMiddleware没有正确运行或请求路径被跳过了认证
	log.Printf("[GetUserIDFromContext] 未找到 parsed_user_id，返回未授权错误")
	return uuid.Nil, &types.AppError{
		Code:    http.StatusUnauthorized,
		Type:    "UNAUTHORIZED",
		Message: "未授权访问4",
	}
}

// HandleUnauthorized 统一处理未授权错误
func HandleUnauthorized(ctx echo.Context, err error) error {
	if appErr, ok := err.(*types.AppError); ok {
		return ctx.JSON(appErr.Code, types.ErrorResponse{
			Success: false,
			Error:   appErr.Type,
			Message: appErr.Message,
		})
	}
	return ctx.JSON(http.StatusUnauthorized, types.ErrorResponse{
		Success: false,
		Error:   "UNAUTHORIZED",
		Message: "未授权访问5",
	})
}
