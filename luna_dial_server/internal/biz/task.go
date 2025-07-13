package biz

import (
	"context"
	"time"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	TaskType    PeriodType `json:"type"`
	TimePeriod  Period     `json:"period"`
	Tags        []string   `json:"tags"`
	Icon        string     `json:"icon"`
	Score       int        `json:"score"`
	IsCompleted bool       `json:"is_completed"`
	ParentID    string     `json:"parent_id"`
	UserID      string     `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// 创建任务参数
type CreateTaskParam struct {
	UserID   string
	Title    string
	Type     PeriodType
	Period   Period
	Tags     []string
	Icon     string
	Score    int
	ParentID string
}

// 编辑任务参数
type UpdateTaskParam struct {
	TaskID      string
	UserID      string
	Title       *string
	Type        *PeriodType
	Period      *Period
	Tags        *[]string
	Icon        *string
	Score       *int
	IsCompleted *bool
}

// 删除任务参数
type DeleteTaskParam struct {
	TaskID string
	UserID string
}

// 设置任务分数参数
type SetTaskScoreParam struct {
	TaskID string
	UserID string
	Score  int
}

// 创建子任务参数
type CreateSubTaskParam struct {
	ParentID string
	UserID   string
	Title    string
	Type     PeriodType
	Period   Period
	Tags     []string
	Icon     string
	Score    int
}

// 添加标签参数
type AddTagParam struct {
	TaskID string
	UserID string
	Tag    string
}

// 移除标签参数
type RemoveTagParam struct {
	TaskID string
	UserID string
	Tag    string
}

// 设置任务icon参数
type SetTaskIconParam struct {
	TaskID string
	UserID string
	Icon   string
}

// 获取指定时间的指定类型的任务列表参数
type ListTaskByPeriodParam struct {
	UserID  string
	Period  Period
	GroupBy PeriodType
}

// 获取某个任务的父任务树列表参数
type ListTaskParentTreeParam struct {
	UserID string
	TaskID string
}

// 获取某个任务的整个任务树参数
type ListTaskTreeParam struct {
	UserID string
	TaskID string
}

type GetTaskStatsParam struct {
	UserID  string
	Period  Period
	GroupBy PeriodType
}

type TaskUsecase struct {
	repo TaskRepo
	// log *log.Helper
}

func NewTaskUsecase(repo TaskRepo) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

func (uc *TaskUsecase) CreateTask(ctx context.Context, param CreateTaskParam) (*Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) UpdateTask(ctx context.Context, param UpdateTaskParam) (*Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) DeleteTask(ctx context.Context, param DeleteTaskParam) error {
	return ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) SetTaskScore(ctx context.Context, param SetTaskScoreParam) (*Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) CreateSubTask(ctx context.Context, param CreateSubTaskParam) (*Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) AddTag(ctx context.Context, param AddTagParam) (*Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) RemoveTag(ctx context.Context, param RemoveTagParam) (*Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) SetTaskIcon(ctx context.Context, param SetTaskIconParam) (*Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) ListTaskByPeriod(ctx context.Context, param ListTaskByPeriodParam) ([]Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) ListTaskTree(ctx context.Context, param ListTaskTreeParam) ([]Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) ListTaskParentTree(ctx context.Context, param ListTaskParentTreeParam) ([]Task, error) {
	return nil, ErrNoPermission // TODO: 实现
}

func (uc *TaskUsecase) GetTaskStats(ctx context.Context, param GetTaskStatsParam) ([]GroupStat, error) {
	return nil, ErrNoPermission // TODO: 实现
}
