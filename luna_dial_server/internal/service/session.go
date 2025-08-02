package service

import (
	"context"
	"errors"
	"luna_dial/internal/data"
	"strings"

	"github.com/labstack/echo/v4"
)

// Session相关的错误定义
var (
	ErrSessionMissing = errors.New("session missing")
)

// SessionMiddleware Session认证中间件
func (s *Service) SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从请求头获取Authorization session
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(401, map[string]string{
					"error":   "unauthorized",
					"message": "Authorization header is required",
				})
			}

			// 验证Bearer session格式
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return echo.NewHTTPError(401, map[string]string{
					"error":   "unauthorized",
					"message": "Authorization header must start with 'Bearer '",
				})
			}

			sessionID := authHeader[len(bearerPrefix):]
			if sessionID == "" {
				return echo.NewHTTPError(401, map[string]string{
					"error":   "unauthorized",
					"message": "Session ID is required",
				})
			}

			// 验证并刷新session
			session, err := s.sessionManager.ValidateSession(c.Request().Context(), sessionID)
			if err != nil {
				var message string
				switch err {
				case data.ErrSessionNotFound:
					message = "Session not found"
				case data.ErrSessionExpired:
					message = "Session has expired"
				case data.ErrSessionInvalid:
					message = "Session is invalid"
				default:
					message = "Session validation failed"
				}

				return echo.NewHTTPError(401, map[string]string{
					"error":   "unauthorized",
					"message": message,
				})
			}

			// 刷新session过期时间
			if err := s.sessionManager.RefreshSession(c.Request().Context(), sessionID); err != nil {
				// 刷新失败时记录日志，但不阻断请求
				// 可以在这里添加日志记录
			}

			// 将用户信息存储到context中
			c.Set("user_id", session.UserID)
			c.Set("username", session.Username)
			c.Set("session_id", session.ID)
			c.Set("session", session)

			return next(c)
		}
	}
}

// OptionalSessionMiddleware 可选的Session认证中间件（不强制要求session）
func (s *Service) OptionalSessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				const bearerPrefix = "Bearer "
				if strings.HasPrefix(authHeader, bearerPrefix) {
					sessionID := authHeader[len(bearerPrefix):]
					if sessionID != "" {
						session, err := s.sessionManager.ValidateSession(c.Request().Context(), sessionID)
						if err == nil {
							// 如果session有效，将用户信息存储到context中
							c.Set("user_id", session.UserID)
							c.Set("username", session.Username)
							c.Set("session_id", session.ID)
							c.Set("session", session)

							// 刷新session过期时间
							s.sessionManager.RefreshSession(c.Request().Context(), sessionID)
						}
					}
				}
			}
			return next(c)
		}
	}
}

// CreateSession 创建新的session
func (s *Service) CreateSession(ctx context.Context, userID string, username string) (*data.SessionResponse, error) {
	return s.sessionManager.CreateSession(ctx, userID, username)
}

// DeleteSession 删除session（用户退出）
func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
	return s.sessionManager.DeleteSession(ctx, sessionID)
}

// DeleteUserSessions 删除用户的所有session
func (s *Service) DeleteUserSessions(ctx context.Context, userID string) error {
	return s.sessionManager.DeleteUserSessions(ctx, userID)
}

// GetSessionFromContext 从Echo Context中获取Session信息
func GetSessionFromContext(c echo.Context) (*data.Session, bool) {
	sessionVal := c.Get("session")
	if sessionVal == nil {
		return nil, false
	}

	session, ok := sessionVal.(*data.Session)
	return session, ok
}

// GetSessionIDFromContext 从Echo Context中获取Session ID
func GetSessionIDFromContext(c echo.Context) (string, bool) {
	sessionIDVal := c.Get("session_id")
	if sessionIDVal == nil {
		return "", false
	}

	sessionID, ok := sessionIDVal.(string)
	return sessionID, ok
}

// 从 ctx 获取用户信息 id，name
func GetUserFromContext(c echo.Context) (string, string, error) {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return "", "", ErrSessionMissing
	}

	userID, ok := userIDVal.(string)
	if !ok {
		return "", "", ErrSessionMissing
	}

	usernameVal := c.Get("username")
	if usernameVal == nil {
		return "", "", ErrSessionMissing
	}

	username, ok := usernameVal.(string)
	if !ok {
		return "", "", ErrSessionMissing
	}

	return username, userID, nil
}
