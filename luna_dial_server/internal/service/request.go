package service

import (
	"fmt"
	"luna_dial/internal/biz"
	"time"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ListTaskRequest struct {
    PeriodType string    `json:"period_type" validate:"required,oneof=day week month quarter year"`
    StartDate  time.Time `json:"start_date" validate:"required"`
    EndDate    time.Time `json:"end_date" validate:"required"`
}

type CreateTaskRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	PeriodType  string    `json:"period_type" validate:"required,oneof=day week month quarter year"`
	Priority    string    `json:"priority" validate:"required,oneof=low medium high urgent"`
	Icon        string    `json:"icon"`
	Tags        []string  `json:"tags"`
}

type CreateSubTaskRequest struct {
    Title       string    `json:"title" validate:"required"`
    Description string    `json:"description"`
    StartDate   time.Time `json:"start_date" validate:"required"`
    EndDate     time.Time `json:"end_date" validate:"required"`
    PeriodType  string    `json:"period_type" validate:"required,oneof=day week month quarter year"`
    Priority    string    `json:"priority" validate:"required,oneof=low medium high urgent"`
    Icon        string    `json:"icon"`
    Tags        []string  `json:"tags"`
    // 兼容旧客户端：允许携带 task_id，但不再校验；服务端使用路径参数作为父任务ID
    TaskID      string    `json:"task_id,omitempty"`
}

// 更新任务
type UpdateTaskRequest struct {
    Title       *string    `json:"title,omitempty"`
    Description *string    `json:"description,omitempty"`
    StartDate   *time.Time `json:"start_date,omitempty"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    Priority    *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high urgent"`
    Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=not_started in_progress completed cancelled"`
    Icon        *string    `json:"icon,omitempty"`
    Tags        *[]string  `json:"tags,omitempty"`
    // 任务ID改由路径参数传入，保留字段以向后兼容
    TaskID      string     `json:"task_id,omitempty"`
}

// 标记任务完成
type CompleteTaskRequest struct {
	TaskID string `json:"task_id" validate:"required"`
}

// 更新任务分数
type UpdateTaskScoreRequest struct {
    Score  int    `json:"score" validate:"required"`
}

// 删除任务
type DeleteTaskRequest struct {
	TaskID string `json:"task_id" validate:"required"`
}

type ListJournalByPeriodRequest struct {
	PeriodType string    `json:"period_type" validate:"required,oneof=day week month quarter year"`
	StartDate  time.Time `json:"start_date" validate:"required"`
	EndDate    time.Time `json:"end_date" validate:"required"`
}

// 新建日志请求
type CreateJournalRequest struct {
	Title       string    `json:"title" validate:"required"`
	Content     string    `json:"content" validate:"required"`
	JournalType string    `json:"journal_type" validate:"required,oneof=day week month quarter year"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`

	Icon string `json:"icon"`
}

// 更新日志请求
type UpdateJournalRequest struct {
    // 日志ID改由路径参数传入
    JournalID   string  `json:"journal_id,omitempty"`
    Title       *string `json:"title,omitempty"`
    Content     *string `json:"content,omitempty"`
    JournalType *string `json:"journal_type,omitempty" validate:"omitempty,oneof=day week month quarter year"`
    // StartDate   *time.Time `json:"start_date,omitempty"`
    // EndDate     *time.Time `json:"end_date,omitempty"`
    Icon *string `json:"icon,omitempty"`
}

// 查看list
type ListPlansRequest struct {
	PeriodType string    `json:"period_type" validate:"required,oneof=day week month quarter year"`
	StartDate  time.Time `json:"start_date" validate:"required"`
	EndDate    time.Time `json:"end_date" validate:"required"`
}

// 分页查询根任务请求
type ListRootTasksRequest struct {
	Page     int      `json:"page" validate:"min=1"`                                                              // 页码，默认1
	PageSize int      `json:"page_size" validate:"min=1,max=100"`                                                 // 每页大小，默认20
	Status   []string `json:"status,omitempty" validate:"dive,oneof=not_started in_progress completed cancelled"` // 状态过滤
	Priority []string `json:"priority,omitempty" validate:"dive,oneof=low medium high urgent"`                    // 优先级过滤
	TaskType []string `json:"task_type,omitempty" validate:"dive,oneof=day week month quarter year"`              // 任务类型过滤
}

// 获取全局任务树请求（分页）
type ListGlobalTaskTreeRequest struct {
	Page         int      `json:"page" validate:"min=1"`                                                              // 页码，默认1
	PageSize     int      `json:"page_size" validate:"min=1,max=50"`                                                  // 每页大小，默认10，最大50
	Status       []string `json:"status,omitempty" validate:"dive,oneof=not_started in_progress completed cancelled"` // 状态过滤
	IncludeEmpty bool     `json:"include_empty,omitempty"`                                                            // 是否包含无子任务的根任务，默认true
}

// 移动任务请求
type MoveTaskRequest struct {
	TaskID      string `json:"task_id" validate:"required"` // 要移动的任务ID
	NewParentID string `json:"new_parent_id,omitempty"`     // 新父任务ID，空表示移动到根级别
}

// 分页查询日志请求（新版本，支持过滤）
type ListJournalsWithPaginationRequest struct {
	Page        int        `json:"page" validate:"min=1"`                                                         // 页码，默认1
	PageSize    int        `json:"page_size" validate:"min=1,max=100"`                                            // 每页大小，默认20
	JournalType *string    `json:"journal_type,omitempty" validate:"omitempty,oneof=day week month quarter year"` // 日志类型过滤
	StartDate   *time.Time `json:"start_date,omitempty"`                                                          // 时间范围过滤开始
	EndDate     *time.Time `json:"end_date,omitempty"`                                                            // 时间范围过滤结束
}

func PeriodTypeFromString(s string) (biz.PeriodType, error) {
	switch s {
	case "day":
		return biz.PeriodDay, nil
	case "week":
		return biz.PeriodWeek, nil
	case "month":
		return biz.PeriodMonth, nil
	case "quarter":
		return biz.PeriodQuarter, nil
	case "year":
		return biz.PeriodYear, nil
	default:
		return 0, fmt.Errorf("unknown period type: %s", s)
	}
}

func TaskStatusFromString(s string) (biz.TaskStatus, error) {
	switch s {
	case "not_started":
		return biz.TaskStatusNotStarted, nil
	case "in_progress":
		return biz.TaskStatusInProgress, nil
	case "completed":
		return biz.TaskStatusCompleted, nil
	case "cancelled":
		return biz.TaskStatusCancelled, nil
	default:
		return 0, fmt.Errorf("unknown task status: %s", s)
	}
}

func TaskPriorityFromString(s string) (biz.TaskPriority, error) {
	switch s {
	case "low":
		return biz.TaskPriorityLow, nil
	case "medium":
		return biz.TaskPriorityMedium, nil
	case "high":
		return biz.TaskPriorityHigh, nil
	case "urgent":
		return biz.TaskPriorityUrgent, nil
	default:
		return 0, fmt.Errorf("unknown task priority: %s", s)
	}
}
