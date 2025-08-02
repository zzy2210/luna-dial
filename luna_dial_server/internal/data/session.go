package data

import (
	"context"
	"errors"
	"time"
)

// Session相关的错误定义
var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionInvalid  = errors.New("session invalid")
)

// Session 会话结构
type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	LastAccessAt time.Time `json:"last_access_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// SessionResponse 会话响应结构
type SessionResponse struct {
	SessionID string `json:"session_id"`
	ExpiresIn int64  `json:"expires_in"`
}

// SessionManager 会话管理器接口
type SessionManager interface {
	// CreateSession 创建新的会话
	CreateSession(ctx context.Context, userID string, username string) (*SessionResponse, error)
	
	// ValidateSession 验证会话是否有效
	ValidateSession(ctx context.Context, sessionID string) (*Session, error)
	
	// RefreshSession 刷新会话过期时间
	RefreshSession(ctx context.Context, sessionID string) error
	
	// DeleteSession 删除会话
	DeleteSession(ctx context.Context, sessionID string) error
	
	// DeleteUserSessions 删除用户的所有会话
	DeleteUserSessions(ctx context.Context, userID string) error
	
	// CleanupExpired 清理过期的会话
	CleanupExpired(ctx context.Context) error
	
	// Close 关闭会话管理器
	Close() error
}

// IsExpired 检查会话是否过期
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// ToResponse 转换为响应结构
func (s *Session) ToResponse() *SessionResponse {
	return &SessionResponse{
		SessionID: s.ID,
		ExpiresIn: s.ExpiresAt.Unix(),
	}
}
