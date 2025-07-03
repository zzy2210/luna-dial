package service

import (
	"context"
	"okr-web/ent"
	"okr-web/internal/types"
	"time"

	"github.com/google/uuid"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (*ent.User, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*ent.User, error)
	GetUserByUsername(ctx context.Context, username string) (*ent.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*ent.User, error)
}

// TaskService 任务服务接口
type TaskService interface {
	CreateTask(ctx context.Context, userID uuid.UUID, req TaskRequest) (*ent.Task, error)
	GetTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*ent.Task, error)
	UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, req TaskRequest) (*ent.Task, error)
	DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error
	GetTaskChildren(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) ([]*ent.Task, error)
	GetTaskTree(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*types.TaskTree, error)
	GetGlobalView(ctx context.Context, userID uuid.UUID) ([]*types.TaskTree, error)
	GetPlanView(ctx context.Context, userID uuid.UUID, req PlanRequest) (*types.PlanResponse, error)
	GetTasksByUser(ctx context.Context, userID uuid.UUID, filters TaskFilters) (*TaskListResponse, error)
	CreateSubTask(ctx context.Context, userID uuid.UUID, parentID uuid.UUID, req TaskRequest) (*ent.Task, error)
	UpdateTaskScore(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, score int) error
	GetContextView(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*ContextView, error)
}

// JournalService 日志服务接口
type JournalService interface {
	CreateJournal(ctx context.Context, userID uuid.UUID, req JournalRequest) (*ent.JournalEntry, error)
	GetJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID) (*ent.JournalEntry, error)
	UpdateJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID, req JournalRequest) (*ent.JournalEntry, error)
	DeleteJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID) error
	GetJournalsByTime(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) ([]*ent.JournalEntry, error)
	GetJournalsByUser(ctx context.Context, userID uuid.UUID, filters JournalFilters) (*JournalListResponse, error)
	LinkJournalToTasks(ctx context.Context, userID uuid.UUID, journalID uuid.UUID, taskIDs []uuid.UUID) error
}

// StatsService 统计服务接口
type StatsService interface {
	GetTaskCompletionStats(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) (*TaskCompletionStats, error)
	GetScoreTrend(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) (*ScoreTrendStats, error)
	GetTimeDistribution(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) (*TimeDistributionStats, error)
	GetUserOverview(ctx context.Context, userID uuid.UUID) (*UserOverviewStats, error)
	GetScoreTrendByReference(ctx context.Context, userID uuid.UUID, req ScoreTrendRequest) (*types.ScoreTrendResponse, error)
}

// 请求结构体定义
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
}

type TaskRequest struct {
	Title       string           `json:"title" validate:"required,min=1,max=200"`
	Description *string          `json:"description,omitempty"`
	Type        types.TaskType   `json:"type" validate:"required"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	Status      types.TaskStatus `json:"status" validate:"required"`
	Score       *int             `json:"score,omitempty" validate:"omitempty,min=0,max=10"`
	Tags        *string          `json:"tags,omitempty"`
}

type JournalRequest struct {
	Content       string          `json:"content" validate:"required"`
	TimeReference string          `json:"time_reference"`
	TimeScale     types.TimeScale `json:"time_scale" validate:"required"`
	EntryType     types.EntryType `json:"entry_type" validate:"required"`
	TaskIDs       []uuid.UUID     `json:"task_ids,omitempty"`
}

type TimeRangeRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type TaskFilters struct {
	Type      *types.TaskType   `json:"type,omitempty"`
	Status    *types.TaskStatus `json:"status,omitempty"`
	ParentID  *uuid.UUID        `json:"parent_id,omitempty"`
	StartDate *time.Time        `json:"start_date,omitempty"`
	EndDate   *time.Time        `json:"end_date,omitempty"`
	Page      int               `json:"page" validate:"min=1"`
	PageSize  int               `json:"page_size" validate:"min=1,max=100"`
}

type JournalFilters struct {
	TimeScale *types.TimeScale `json:"time_scale,omitempty"`
	EntryType *types.EntryType `json:"entry_type,omitempty"`
	StartDate *time.Time       `json:"start_date,omitempty"`
	EndDate   *time.Time       `json:"end_date,omitempty"`
	Page      int              `json:"page" validate:"min=1"`
	PageSize  int              `json:"page_size" validate:"min=1,max=100"`
}

// PlanRequest 计划视图请求结构体
type PlanRequest struct {
	Scale   types.TimeScale `json:"scale" query:"scale" validate:"required"`       // 时间尺度
	TimeRef string          `json:"time_ref" query:"time_ref" validate:"required"` // 时间参考
}

// ScoreTrendRequest 分数趋势请求结构体
type ScoreTrendRequest struct {
	Scale   types.TimeScale `json:"scale" query:"scale" validate:"required"`       // 统计尺度
	TimeRef string          `json:"time_ref" query:"time_ref" validate:"required"` // 时间参考
}

// 响应结构体定义
type AuthResponse struct {
	User  *ent.User `json:"user"`
	Token string    `json:"token"`
}

type ContextView struct {
	Current  *ent.Task   `json:"current"`
	Parent   *ent.Task   `json:"parent,omitempty"`
	Children []*ent.Task `json:"children"`
	Siblings []*ent.Task `json:"siblings"`
}

type TaskListResponse struct {
	Tasks       []*ent.Task `json:"tasks"`
	Total       int64       `json:"total"`
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
	TotalPages  int         `json:"total_pages"`
}

type JournalListResponse struct {
	Journals    []*ent.JournalEntry `json:"journals"`
	Total       int64               `json:"total"`
	CurrentPage int                 `json:"current_page"`
	PageSize    int                 `json:"page_size"`
	TotalPages  int                 `json:"total_pages"`
}

type TaskCompletionStats struct {
	TotalTasks     int64                      `json:"total_tasks"`
	CompletedTasks int64                      `json:"completed_tasks"`
	CompletionRate float64                    `json:"completion_rate"`
	ByType         map[types.TaskType]int64   `json:"by_type"`
	ByStatus       map[types.TaskStatus]int64 `json:"by_status"`
}

type ScoreTrendStats struct {
	AverageScore float64           `json:"average_score"`
	Trend        []ScoreTrendPoint `json:"trend"`
}

type ScoreTrendPoint struct {
	Date  time.Time `json:"date"`
	Score float64   `json:"score"`
}

type TimeDistributionStats struct {
	ByType      map[types.TaskType]int64  `json:"by_type"`
	ByTimeScale map[types.TimeScale]int64 `json:"by_time_scale"`
}

type UserOverviewStats struct {
	TotalTasks       int64   `json:"total_tasks"`
	CompletedTasks   int64   `json:"completed_tasks"`
	AverageScore     float64 `json:"average_score"`
	TotalJournals    int64   `json:"total_journals"`
	ActiveTasksCount int64   `json:"active_tasks_count"`
}
