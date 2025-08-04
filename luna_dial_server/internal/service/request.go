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
	PeriodType string    `json:"period_type" validate:"required,oneof=daily weekly monthly yearly"`
	StartDate  time.Time `json:"start_date" validate:"required"`
	EndDate    time.Time `json:"end_date" validate:"required"`
}

type CreateTaskRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Priority    string    `json:"priority" validate:"required,oneof=low medium high"`
	Icon        string    `json:"icon"`
	Tags        []string  `json:"tags"`
}

type CreateSubTaskRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Priority    string    `json:"priority" validate:"required,oneof=low medium high"`
	Icon        string    `json:"icon"`
	Tags        []string  `json:"tags"`
	TaskID      string    `json:"task_id" validate:"required"`
}

// 更新任务
type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Priority    *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	Icon        *string    `json:"icon,omitempty"`
	Tags        *[]string  `json:"tags,omitempty"`
	TaskID      string     `json:"task_id" validate:"required"`
	IsCompleted *bool      `json:"is_completed,omitempty"`
}

// 标记任务完成
type CompleteTaskRequest struct {
	TaskID string `json:"task_id" validate:"required"`
}

// 更新任务分数
type UpdateTaskScoreRequest struct {
	TaskID string `json:"task_id" validate:"required"`
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
	JournalID   string  `json:"journal_id" validate:"required"`
	Title       *string `json:"title,omitempty"`
	Content     *string `json:"content,omitempty"`
	JournalType *string `json:"journal_type,omitempty" validate:"omitempty,oneof=day week month quarter year"`
	// StartDate   *time.Time `json:"start_date,omitempty"`
	// EndDate     *time.Time `json:"end_date,omitempty"`
	Icon *string `json:"icon,omitempty"`
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
