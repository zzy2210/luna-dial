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
}
