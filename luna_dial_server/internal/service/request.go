package service

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Response 通用响应结构体
type Response struct {
	Code      int         `json:"code"`                 // 业务状态码
	Message   string      `json:"message"`              // 响应消息
	Data      interface{} `json:"data,omitempty"`       // 响应数据
	Success   bool        `json:"success"`              // 请求是否成功
	Timestamp int64       `json:"timestamp"`            // 响应时间戳
	RequestID string      `json:"request_id,omitempty"` // 请求ID，用于追踪
}

// 登录响应模型
type LoginResponse struct {
	SessionID string `json:"session_id"` // 会话ID
	ExpiresIn int64  `json:"expires_in"` // 会话过期时间（秒）
}
