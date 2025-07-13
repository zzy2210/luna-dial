package biz

import "time"

type TaskRepo interface {
	CreateTask(task *Task) error
	UpdateTask(task *Task) error
	DeleteTask(taskID, userID string) error
	GetTask(taskID, userID string) (*Task, error)
	ListTasks(userID string, periodStart, periodEnd time.Time, taskType string) ([]*Task, error)
	ListTaskTree(taskID, userID string) ([]*Task, error)
}
