package data

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// MemorySessionManager 内存会话管理器
type MemorySessionManager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	
	// 配置
	sessionTimeout  time.Duration
	cleanupInterval time.Duration
	
	// 清理协程控制
	stopCleanup chan bool
	cleanupDone chan bool
}

// NewMemorySessionManager 创建内存会话管理器
func NewMemorySessionManager(sessionTimeout time.Duration) *MemorySessionManager {
	cleanupInterval := 10 * time.Minute // 每10分钟清理一次
	
	manager := &MemorySessionManager{
		sessions:        make(map[string]*Session),
		sessionTimeout:  sessionTimeout,
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan bool),
		cleanupDone:     make(chan bool),
	}
	
	// 启动清理协程
	go manager.cleanupWorker()
	
	return manager
}

// generateSessionID 生成会话ID
func (m *MemorySessionManager) generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256位随机数
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateSession 创建新的会话
func (m *MemorySessionManager) CreateSession(ctx context.Context, userID int64, username string) (*SessionResponse, error) {
	sessionID, err := m.generateSessionID()
	if err != nil {
		return nil, err
	}
	
	now := time.Now()
	session := &Session{
		ID:           sessionID,
		UserID:       userID,
		Username:     username,
		CreatedAt:    now,
		LastAccessAt: now,
		ExpiresAt:    now.Add(m.sessionTimeout),
	}
	
	m.mutex.Lock()
	m.sessions[sessionID] = session
	m.mutex.Unlock()
	
	return session.ToResponse(), nil
}

// ValidateSession 验证会话是否有效
func (m *MemorySessionManager) ValidateSession(ctx context.Context, sessionID string) (*Session, error) {
	m.mutex.RLock()
	session, exists := m.sessions[sessionID]
	m.mutex.RUnlock()
	
	if !exists {
		return nil, ErrSessionNotFound
	}
	
	if session.IsExpired() {
		// 删除过期会话
		m.mutex.Lock()
		delete(m.sessions, sessionID)
		m.mutex.Unlock()
		return nil, ErrSessionExpired
	}
	
	return session, nil
}

// RefreshSession 刷新会话过期时间
func (m *MemorySessionManager) RefreshSession(ctx context.Context, sessionID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	session, exists := m.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}
	
	if session.IsExpired() {
		delete(m.sessions, sessionID)
		return ErrSessionExpired
	}
	
	// 刷新过期时间和最后访问时间
	now := time.Now()
	session.LastAccessAt = now
	session.ExpiresAt = now.Add(m.sessionTimeout)
	
	return nil
}

// DeleteSession 删除会话
func (m *MemorySessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	delete(m.sessions, sessionID)
	return nil
}

// DeleteUserSessions 删除用户的所有会话
func (m *MemorySessionManager) DeleteUserSessions(ctx context.Context, userID int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	for sessionID, session := range m.sessions {
		if session.UserID == userID {
			delete(m.sessions, sessionID)
		}
	}
	
	return nil
}

// CleanupExpired 清理过期的会话
func (m *MemorySessionManager) CleanupExpired(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	now := time.Now()
	expiredSessions := make([]string, 0)
	
	for sessionID, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			expiredSessions = append(expiredSessions, sessionID)
		}
	}
	
	for _, sessionID := range expiredSessions {
		delete(m.sessions, sessionID)
	}
	
	return nil
}

// cleanupWorker 清理工作协程
func (m *MemorySessionManager) cleanupWorker() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			m.CleanupExpired(context.Background())
		case <-m.stopCleanup:
			m.cleanupDone <- true
			return
		}
	}
}

// Close 关闭会话管理器
func (m *MemorySessionManager) Close() error {
	// 停止清理协程
	m.stopCleanup <- true
	<-m.cleanupDone
	
	// 清理所有会话
	m.mutex.Lock()
	m.sessions = make(map[string]*Session)
	m.mutex.Unlock()
	
	return nil
}

// GetSessionCount 获取当前会话数量（用于调试）
func (m *MemorySessionManager) GetSessionCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.sessions)
}

// GetUserSessions 获取用户的所有会话（用于调试）
func (m *MemorySessionManager) GetUserSessions(userID int64) []*Session {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	sessions := make([]*Session, 0)
	for _, session := range m.sessions {
		if session.UserID == userID && !session.IsExpired() {
			sessions = append(sessions, session)
		}
	}
	
	return sessions
}
