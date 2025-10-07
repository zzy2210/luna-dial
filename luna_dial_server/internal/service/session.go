package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"luna_dial/internal/data"
)

// Session相关的错误定义
var (
	ErrSessionMissing = errors.New("session missing")
)

// SessionMiddleware Session认证中间件
func (s *Service) SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var sessionID string

			// 首先尝试从Authorization header获取session
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				const bearerPrefix = "Bearer "
				if strings.HasPrefix(authHeader, bearerPrefix) {
					sessionID = authHeader[len(bearerPrefix):]
				}
			}

			// 如果Authorization header中没有，则尝试从Cookie获取
			if sessionID == "" {
				cookie, err := c.Cookie("session_id")
				if err == nil {
					sessionID = cookie.Value
				}
			}

			// 如果都没有找到session ID，返回统一错误响应
			if sessionID == "" {
				return c.JSON(401, &Response{
					Code:      401,
					Message:   "Authorization header or session cookie is required",
					Success:   false,
					Timestamp: time.Now().Unix(),
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

				return c.JSON(401, &Response{
					Code:      401,
					Message:   message,
					Success:   false,
					Timestamp: time.Now().Unix(),
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

	return userID, username, nil
}
