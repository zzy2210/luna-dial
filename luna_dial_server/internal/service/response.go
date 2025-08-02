package service

import "time"

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

// PaginatedData 通用分页数据结构 (嵌套在 Response.Data 中)
type PaginatedData struct {
	Items      interface{} `json:"items"`      // 数据列表
	Pagination *Pagination `json:"pagination"` // 分页信息
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页大小
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"total_pages"` // 总页数
	HasNext    bool  `json:"has_next"`    // 是否有下一页
	HasPrev    bool  `json:"has_prev"`    // 是否有上一页
}

// 辅助函数

// NewPagination 创建分页信息
func NewPagination(page, pageSize int, total int64) *Pagination {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	
	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Code:      200,
		Message:   "success",
		Data:      data,
		Success:   true,
		Timestamp: time.Now().Unix(),
	}
}

// NewSuccessResponseWithMessage 创建带自定义消息的成功响应
func NewSuccessResponseWithMessage(message string, data interface{}) *Response {
	return &Response{
		Code:      200,
		Message:   message,
		Data:      data,
		Success:   true,
		Timestamp: time.Now().Unix(),
	}
}

// NewPaginatedResponse 创建分页响应
func NewPaginatedResponse(items interface{}, page, pageSize int, total int64) *Response {
	pagination := NewPagination(page, pageSize, total)
	
	return &Response{
		Code:    200,
		Message: "success",
		Data: &PaginatedData{
			Items:      items,
			Pagination: pagination,
		},
		Success:   true,
		Timestamp: time.Now().Unix(),
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) *Response {
	return &Response{
		Code:      code,
		Message:   message,
		Success:   false,
		Timestamp: time.Now().Unix(),
	}
}

// WithRequestID 添加请求ID
func (r *Response) WithRequestID(requestID string) *Response {
	r.RequestID = requestID
	return r
}
