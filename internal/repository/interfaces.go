package repository

import (
	"context"
	"time"

	"okr-web/ent"
	"okr-web/ent/journalentry"

	"github.com/google/uuid"
)

// ScorePoint 表示分数统计点
type ScorePoint struct {
	Date  time.Time `json:"date"`
	Score int       `json:"score"`
	Count int       `json:"count"`
}

// UserRepository 用户Repository接口
type UserRepository interface {
	Create(ctx context.Context, builder func(*ent.UserCreate) *ent.UserCreate) (*ent.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error)
	GetByUsername(ctx context.Context, username string) (*ent.User, error)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	Update(ctx context.Context, id uuid.UUID, updater func(*ent.UserUpdateOne) *ent.UserUpdateOne) (*ent.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// TaskRepository 任务Repository接口
type TaskRepository interface {
	Create(ctx context.Context, builder func(*ent.TaskCreate) *ent.TaskCreate) (*ent.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Task, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ent.Task, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*ent.Task, error)
	GetRootTasks(ctx context.Context, userID uuid.UUID) ([]*ent.Task, error)
	GetTaskTree(ctx context.Context, userID uuid.UUID, rootID *uuid.UUID) ([]*ent.Task, error)
	Update(ctx context.Context, id uuid.UUID, updater func(*ent.TaskUpdateOne) *ent.TaskUpdateOne) (*ent.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	GetTasksByTimeRange(ctx context.Context, userID uuid.UUID, start, end time.Time, taskType string) ([]*ent.Task, error)
	GetTasksWithAncestors(ctx context.Context, userID uuid.UUID, start, end time.Time, taskType string) ([]*ent.Task, error)
	GetScoreStatsInTimeRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]ScorePoint, error)
}

// JournalRepository 日志Repository接口
type JournalRepository interface {
	Create(ctx context.Context, builder func(*ent.JournalEntryCreate) *ent.JournalEntryCreate) (*ent.JournalEntry, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ent.JournalEntry, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ent.JournalEntry, error)
	GetByTimeReference(ctx context.Context, userID uuid.UUID, timeRef string) ([]*ent.JournalEntry, error)
	Update(ctx context.Context, id uuid.UUID, updater func(*ent.JournalEntryUpdateOne) *ent.JournalEntryUpdateOne) (*ent.JournalEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	GetByTimeRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*ent.JournalEntry, error)
	// 新增：按时间尺度和 time_reference 批量获取日志
	GetByTimeScaleAndReferences(ctx context.Context, userID uuid.UUID, timeScale journalentry.TimeScale, timeRefs []string) ([]*ent.JournalEntry, error)
}
