package repository

import (
	"context"
	"fmt"
	"time"

	"okr-web/ent"
	"okr-web/ent/task"

	"github.com/google/uuid"
)

// taskRepository 任务Repository实现
type taskRepository struct {
	client *ent.Client
}

// NewTaskRepository 创建新的任务Repository
func NewTaskRepository(client *ent.Client) TaskRepository {
	return &taskRepository{client: client}
}

// Create 创建新任务
func (r *taskRepository) Create(ctx context.Context, builder func(*ent.TaskCreate) *ent.TaskCreate) (*ent.Task, error) {
	return builder(r.client.Task.Create()).Save(ctx)
}

// GetByID 根据ID获取任务
func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Task, error) {
	t, err := r.client.Task.
		Query().
		Where(task.ID(id)).
		WithUser().
		WithParent().
		WithChildren().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return t, nil
}

// GetByUserID 根据用户ID获取任务列表
func (r *taskRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ent.Task, error) {
	query := r.client.Task.
		Query().
		Where(task.UserID(userID)).
		WithUser().
		WithParent().
		Order(ent.Desc(task.FieldCreatedAt))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	tasks, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by user ID: %w", err)
	}
	return tasks, nil
}

// GetChildren 获取子任务
func (r *taskRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*ent.Task, error) {
	children, err := r.client.Task.
		Query().
		Where(task.ParentID(parentID)).
		WithUser().
		Order(ent.Asc(task.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get children tasks: %w", err)
	}
	return children, nil
}

// GetRootTasks 获取根任务（没有父任务的任务）
func (r *taskRepository) GetRootTasks(ctx context.Context, userID uuid.UUID) ([]*ent.Task, error) {
	tasks, err := r.client.Task.
		Query().
		Where(
			task.UserID(userID),
			task.ParentIDIsNil(),
		).
		WithUser().
		Order(ent.Asc(task.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get root tasks: %w", err)
	}
	return tasks, nil
}

// GetTaskTree 获取任务树（递归获取所有子任务）
func (r *taskRepository) GetTaskTree(ctx context.Context, userID uuid.UUID, rootID *uuid.UUID) ([]*ent.Task, error) {
	var query *ent.TaskQuery

	if rootID != nil {
		// 获取指定根任务及其所有子任务
		query = r.client.Task.
			Query().
			Where(
				task.UserID(userID),
				task.Or(
					task.ID(*rootID),
					task.HasParentWith(task.ID(*rootID)),
				),
			)
	} else {
		// 获取用户的所有任务
		query = r.client.Task.
			Query().
			Where(task.UserID(userID))
	}

	tasks, err := query.
		WithUser().
		WithParent().
		WithChildren().
		Order(ent.Asc(task.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get task tree: %w", err)
	}

	return tasks, nil
}

// Update 更新任务
func (r *taskRepository) Update(ctx context.Context, id uuid.UUID, updater func(*ent.TaskUpdateOne) *ent.TaskUpdateOne) (*ent.Task, error) {
	updateOne := r.client.Task.UpdateOneID(id)
	updateOne = updater(updateOne)

	t, err := updateOne.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return t, nil
}

// Delete 删除任务
func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.Task.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("task not found")
		}
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// CountByUserID 统计用户的任务数量
func (r *taskRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := r.client.Task.
		Query().
		Where(task.UserID(userID)).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %w", err)
	}
	return count, nil
}

// GetTasksByTimeRange 根据时间范围和类型获取任务
func (r *taskRepository) GetTasksByTimeRange(ctx context.Context, userID uuid.UUID, start, end time.Time, taskType string) ([]*ent.Task, error) {
	tasks, err := r.client.Task.
		Query().
		Where(
			task.UserID(userID),
			task.TypeEQ(task.Type(taskType)),
			task.And(
				task.Or(
					task.StartDateGTE(start),
					task.EndDateGTE(start),
					task.CreatedAtGTE(start),
				),
				task.Or(
					task.StartDateLTE(end),
					task.EndDateLTE(end),
					task.CreatedAtLTE(end),
				),
			),
		).
		WithUser().
		WithParent().
		WithChildren().
		Order(ent.Asc(task.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by time range: %w", err)
	}
	return tasks, nil
}

// GetTasksWithAncestors 获取时间范围内的任务及其完整父级链（支持类型过滤）
func (r *taskRepository) GetTasksWithAncestors(ctx context.Context, userID uuid.UUID, start, end time.Time, taskType string) ([]*ent.Task, error) {
	// 首先获取时间范围内的任务（带类型过滤）
	directTasks, err := r.GetTasksByTimeRange(ctx, userID, start, end, taskType)
	if err != nil {
		return nil, err
	}

	// 收集所有需要的任务ID（包括父级链）
	allTaskIDs := make(map[uuid.UUID]struct{})
	var collectAncestors func(taskID uuid.UUID) error

	collectAncestors = func(taskID uuid.UUID) error {
		if _, exists := allTaskIDs[taskID]; exists {
			return nil // 已经处理过
		}

		allTaskIDs[taskID] = struct{}{}

		// 获取父任务
		t, err := r.client.Task.
			Query().
			Where(task.ID(taskID)).
			WithParent().
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil // 任务不存在，跳过
			}
			return err
		}

		// 递归处理父任务
		if t.Edges.Parent != nil {
			return collectAncestors(t.Edges.Parent.ID)
		}

		return nil
	}

	// 收集所有直接任务的父级链
	for _, task := range directTasks {
		if err := collectAncestors(task.ID); err != nil {
			return nil, fmt.Errorf("failed to collect ancestors: %w", err)
		}
	}

	// 获取所有相关任务
	taskIDs := make([]uuid.UUID, 0, len(allTaskIDs))
	for taskID := range allTaskIDs {
		taskIDs = append(taskIDs, taskID)
	}

	allTasks, err := r.client.Task.
		Query().
		Where(
			task.UserID(userID),
			task.IDIn(taskIDs...),
		).
		WithUser().
		WithParent().
		WithChildren().
		Order(ent.Asc(task.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks with ancestors: %w", err)
	}

	return allTasks, nil
}

// GetScoreStatsInTimeRange 获取时间范围内的分数统计
func (r *taskRepository) GetScoreStatsInTimeRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]ScorePoint, error) {
	// 获取时间范围内的任务（统计时不过滤类型，兼容原有逻辑，传空字符串）
	tasks, err := r.GetTasksByTimeRange(ctx, userID, start, end, "")
	if err != nil {
		return nil, err
	}

	// 按日期分组统计
	scoreMap := make(map[string]*ScorePoint)

	for _, t := range tasks {
		var date time.Time

		// 使用任务的结束日期
		date = t.EndDate

		dateStr := date.Format("2006-01-02")

		if point, exists := scoreMap[dateStr]; exists {
			point.Score += t.Score
			point.Count++
		} else {
			scoreMap[dateStr] = &ScorePoint{
				Date:  date,
				Score: t.Score,
				Count: 1,
			}
		}
	}

	// 转换为切片
	var result []ScorePoint
	for _, point := range scoreMap {
		result = append(result, *point)
	}

	return result, nil
}
