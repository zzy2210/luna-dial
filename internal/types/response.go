package types

import (
	"net/http"

	"github.com/google/uuid"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse 统一成功响应格式
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// PaginationResponse 分页响应格式
type PaginationResponse struct {
	Success     bool        `json:"success"`
	Data        interface{} `json:"data"`
	Total       int64       `json:"total"`
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
	TotalPages  int         `json:"total_pages"`
}

// AppError 应用错误类型
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *AppError) Error() string {
	return e.Message
}

// 预定义错误类型
var (
	ErrUserNotFound    = &AppError{Code: http.StatusNotFound, Message: "用户不存在", Type: "USER_NOT_FOUND"}
	ErrTaskNotFound    = &AppError{Code: http.StatusNotFound, Message: "任务不存在", Type: "TASK_NOT_FOUND"}
	ErrJournalNotFound = &AppError{Code: http.StatusNotFound, Message: "日志不存在", Type: "JOURNAL_NOT_FOUND"}
	ErrUnauthorized    = &AppError{Code: http.StatusUnauthorized, Message: "未授权访问3", Type: "UNAUTHORIZED"}
	ErrInvalidRequest  = &AppError{Code: http.StatusBadRequest, Message: "请求参数无效", Type: "INVALID_REQUEST"}
	ErrInternalServer  = &AppError{Code: http.StatusInternalServerError, Message: "服务器内部错误", Type: "INTERNAL_SERVER_ERROR"}
	ErrDuplicateUser   = &AppError{Code: http.StatusConflict, Message: "用户名或邮箱已存在", Type: "DUPLICATE_USER"}
	ErrInvalidPassword = &AppError{Code: http.StatusUnauthorized, Message: "密码错误", Type: "INVALID_PASSWORD"}
)

// 计划视图相关结构体

// PlanRequest 定义了获取计划视图的请求参数
type PlanRequest struct {
	Scale   TimeScale `query:"scale" validate:"required" json:"scale"`       // 时间尺度 (e.g., "quarter", "month")
	TimeRef string    `query:"time_ref" validate:"required" json:"time_ref"` // 时间参考 (e.g., "2024-Q4", "2025-07")
}

// TaskTree 表示包含完整父级链的任务树结构
type TaskTree struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	StartDate   string     `json:"start_date"`
	EndDate     string     `json:"end_date"`
	Status      string     `json:"status"`
	Score       int        `json:"score"`
	ParentID    *uuid.UUID `json:"parent_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Tags        string     `json:"tags"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	// 完整的父级链，从根任务到直接父任务
	Ancestors []*TaskTree `json:"ancestors,omitempty"`
	// 直接子任务
	Children []*TaskTree `json:"children,omitempty"`
	// 任务层级深度（根任务为0）
	Depth int `json:"depth"`
}

// PlanResponse 定义了计划视图的响应数据
type PlanResponse struct {
	// 包含完整父级链的任务树列表
	Tasks []*TaskTree `json:"tasks"`
	// 相关时间周期的日志列表
	Journals interface{} `json:"journals"`
	// 时间范围信息
	TimeRange *TimeRange `json:"time_range"`
	// 统计信息
	Stats *PlanStats `json:"stats"`
}

// PlanStats 计划视图统计信息
type PlanStats struct {
	TotalTasks      int `json:"total_tasks"`       // 总任务数
	CompletedTasks  int `json:"completed_tasks"`   // 已完成任务数
	InProgressTasks int `json:"in_progress_tasks"` // 进行中任务数
	PendingTasks    int `json:"pending_tasks"`     // 待开始任务数
	TotalScore      int `json:"total_score"`       // 总分
	CompletedScore  int `json:"completed_score"`   // 已完成任务总分
}

// 分数趋势统计相关结构体

// ScoreTrendRequest 定义了获取分数趋势的请求参数
type ScoreTrendRequest struct {
	Scale   TimeScale `query:"scale" validate:"required" json:"scale"`       // 统计尺度 (e.g., "month", "quarter")
	TimeRef string    `query:"time_ref" validate:"required" json:"time_ref"` // 时间参考 (e.g., "2025-07", "2024-Q4")
}

// ScoreTrendResponse 定义了分数趋势的响应数据
type ScoreTrendResponse struct {
	Labels []string `json:"labels"` // 时间标签 (e.g., ["2025-07-01", "2025-07-02"])
	Scores []int    `json:"scores"` // 对应标签的分数总和
	Counts []int    `json:"counts"` // 对应标签的任务数量
	// 统计元信息
	Scale     TimeScale     `json:"scale"`      // 统计尺度
	TimeRef   string        `json:"time_ref"`   // 时间参考
	TimeRange *TimeRange    `json:"time_range"` // 实际时间范围
	Summary   *TrendSummary `json:"summary"`    // 趋势摘要
}

// TrendSummary 趋势统计摘要
type TrendSummary struct {
	TotalScore       int     `json:"total_score"`        // 总分
	TotalTasks       int     `json:"total_tasks"`        // 总任务数
	AverageScore     float64 `json:"average_score"`      // 平均分
	AverageTaskCount float64 `json:"average_task_count"` // 平均任务数
	MaxScore         int     `json:"max_score"`          // 最高分
	MaxTasks         int     `json:"max_tasks"`          // 最多任务数
	MinScore         int     `json:"min_score"`          // 最低分
	MinTasks         int     `json:"min_tasks"`          // 最少任务数
}
