package biz

import (
	"context"
	"time"
)

type TaskRepo interface {
	CreateTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, taskID, userID string) error
	GetTask(ctx context.Context, taskID, userID string) (*Task, error)
	ListTasks(ctx context.Context, userID string, periodStart, periodEnd time.Time, taskType int) ([]*Task, error)
	ListTaskParentTree(ctx context.Context, taskID, userID string) ([]*Task, error)
	ListRootTasksWithPagination(ctx context.Context, userID string, page, pageSize int, includeStatus []TaskStatus) ([]*Task, int64, error)
	ListTasksByRootIDs(ctx context.Context, userID string, rootTaskIDs []string, includeStatus []TaskStatus) ([]*Task, error)
	GetCompleteTaskTree(ctx context.Context, taskID, userID string, includeStatus []TaskStatus) ([]*Task, error)
	GetTaskParentChain(ctx context.Context, taskID, userID string) ([]*Task, error)
	UpdateTreeOptimizationFields(ctx context.Context, taskID, userID string) error
}
