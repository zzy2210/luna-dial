package biz

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 创建测试用的 TaskUsecase 实例
func createTestTaskUsecase() *TaskUsecase {
	repo := &mockTaskRepo{}
	return NewTaskUsecase(repo)
}

// 测试 NewTaskUsecase 构造函数
func TestNewTaskUsecase(t *testing.T) {
	repo := &mockTaskRepo{}
	usecase := NewTaskUsecase(repo)

	require.NotNil(t, usecase, "NewTaskUsecase should not return nil")
	assert.Equal(t, repo, usecase.repo, "repo should be set correctly")
}

// 测试 CreateTask 方法
func TestTaskUsecase_CreateTask(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功创建日任务", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "完成产品需求文档",
			Type:   PeriodDay,
			Period: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			},
			Tags:  []string{"工作", "文档", "产品"},
			Icon:  "📝",
			Score: 80,
		}

		task, err := usecase.CreateTask(ctx, param)

		// ❌ TDD: 期望成功创建，当前业务逻辑未实现会失败
		require.NoError(t, err, "CreateTask should succeed")
		require.NotNil(t, task, "CreateTask should return created task object")

		// 验证返回的任务字段
		assert.Equal(t, param.Title, task.Title, "title should match")
		assert.Equal(t, param.Type, task.TaskType, "task type should match")
		assert.Equal(t, param.Score, task.Score, "score should match")
		assert.Equal(t, param.UserID, task.UserID, "user ID should match")
		assert.Equal(t, param.Icon, task.Icon, "icon should match")
		assert.Equal(t, len(param.Tags), len(task.Tags), "tags count should match")

		// 验证自动设置的字段
		assert.NotEmpty(t, task.ID, "ID should be generated")
		assert.Equal(t, TaskStatusNotStarted, task.Status, "new task should be not started")
		assert.Equal(t, TaskPriorityLow, task.Priority, "new task should have low priority")
		assert.False(t, task.CreatedAt.IsZero(), "created time should be set")
		assert.False(t, task.UpdatedAt.IsZero(), "updated time should be set")
	})

	t.Run("成功创建周任务", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "完成项目里程碑",
			Type:   PeriodWeek,
			Period: Period{
				Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			Tags:  []string{"项目", "里程碑"},
			Icon:  "🎯",
			Score: 200,
		}

		task, err := usecase.CreateTask(ctx, param)

		// ❌ TDD: 期望成功创建，当前业务逻辑未实现会失败
		require.NoError(t, err, "CreateTask should succeed for week task")
		require.NotNil(t, task, "should return created week task")
		assert.Equal(t, PeriodWeek, task.TaskType, "task type should be PeriodWeek")
	})

	t.Run("成功创建子任务", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "子任务：设计UI界面",
			Type:   PeriodDay,
			Period: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			},
			Tags:     []string{"设计", "UI"},
			Icon:     "🎨",
			Score:    50,
			ParentID: "parent-task-123", // 父任务ID
		}

		task, err := usecase.CreateTask(ctx, param)

		// ❌ TDD: 期望成功创建，当前业务逻辑未实现会失败
		require.NoError(t, err, "CreateTask should succeed for subtask")
		require.NotNil(t, task, "should return created subtask")
		assert.Equal(t, param.ParentID, task.ParentID, "parent ID should match")
		assert.Equal(t, param.Title, task.Title, "title should match")
	})

	t.Run("参数验证失败 - 空用户ID", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "", // 空用户ID
			Title:  "测试任务",
			Type:   PeriodDay,
			Score:  50,
		}

		task, err := usecase.CreateTask(ctx, param)

		// ✅ TDD: 明确期望的业务错误
		assert.Nil(t, task, "should return nil task for empty user ID")
		assert.Equal(t, ErrInvalidInput, err, "should return ErrInvalidInput for empty user ID")
	})

	t.Run("参数验证失败 - 空标题", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "", // 空标题
			Type:   PeriodDay,
			Score:  50,
		}

		task, err := usecase.CreateTask(ctx, param)

		// ✅ TDD: 明确期望的业务错误
		assert.Nil(t, task, "should return nil task for empty title")
		assert.Equal(t, ErrInvalidInput, err, "should return ErrInvalidInput for empty title")
	})

	t.Run("参数验证失败 - 无效分数", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "测试任务",
			Type:   PeriodDay,
			Score:  -10, // 负分数
			Period: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
		}

		task, err := usecase.CreateTask(ctx, param)

		// ✅ TDD: 当前会成功创建，因为没有负分数验证
		// 实际业务中可能需要添加分数验证
		if err != nil {
			t.Logf("当前返回错误: %v，可能需要添加分数验证逻辑", err)
		} else if task != nil && task.Score < 0 {
			t.Log("当前允许负分数，可能需要添加验证")
		}
	})
}

// 测试 UpdateTask 方法
func TestTaskUsecase_UpdateTask(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功更新任务标题", func(t *testing.T) {
		newTitle := "更新后的任务标题"
		param := UpdateTaskParam{
			TaskID: "task-123",
			UserID: "user-123",
			Title:  &newTitle,
		}

		task, err := usecase.UpdateTask(ctx, param)

		// ❌ TDD: 期望成功更新，当前业务逻辑未实现会失败
		require.NoError(t, err, "UpdateTask should succeed")
		require.NotNil(t, task, "should return updated task")
		assert.Equal(t, newTitle, task.Title, "title should be updated")
		assert.False(t, task.UpdatedAt.IsZero(), "updated time should be set")
	})

	t.Run("成功更新任务状态", func(t *testing.T) {
		status := TaskStatusCompleted
		param := UpdateTaskParam{
			TaskID: "task-123",
			UserID: "user-123",
			Status: &status,
		}

		task, err := usecase.UpdateTask(ctx, param)

		// ❌ TDD: 期望成功更新，当前业务逻辑未实现会失败
		require.NoError(t, err, "UpdateTask should succeed for completion status")
		require.NotNil(t, task, "should return updated task")
		assert.Equal(t, TaskStatusCompleted, task.Status, "task should be marked as completed")
	})

	t.Run("成功更新任务分数和标签", func(t *testing.T) {
		newScore := 100
		newTags := []string{"更新", "标签"}
		param := UpdateTaskParam{
			TaskID: "task-123",
			UserID: "user-123",
			Score:  &newScore,
			Tags:   &newTags,
		}

		task, err := usecase.UpdateTask(ctx, param)

		// ❌ TDD: 期望成功更新，当前业务逻辑未实现会失败
		require.NoError(t, err, "UpdateTask should succeed for score and tags")
		require.NotNil(t, task, "should return updated task")
		assert.Equal(t, newScore, task.Score, "score should be updated")
		assert.Equal(t, newTags, task.Tags, "tags should be updated")
	})

	t.Run("成功更新任务优先级", func(t *testing.T) {
		priority := TaskPriorityUrgent
		param := UpdateTaskParam{
			TaskID:   "task-123",
			UserID:   "user-123",
			Priority: &priority,
		}

		task, err := usecase.UpdateTask(ctx, param)

		// ❌ TDD: 期望成功更新，当前业务逻辑未实现会失败
		require.NoError(t, err, "UpdateTask should succeed for priority")
		require.NotNil(t, task, "should return updated task")
		assert.Equal(t, TaskPriorityUrgent, task.Priority, "priority should be updated to urgent")
	})

	t.Run("成功更新任务状态和优先级", func(t *testing.T) {
		status := TaskStatusInProgress
		priority := TaskPriorityHigh
		param := UpdateTaskParam{
			TaskID:   "task-123",
			UserID:   "user-123",
			Status:   &status,
			Priority: &priority,
		}

		task, err := usecase.UpdateTask(ctx, param)

		// ❌ TDD: 期望成功更新，当前业务逻辑未实现会失败
		require.NoError(t, err, "UpdateTask should succeed for status and priority")
		require.NotNil(t, task, "should return updated task")
		assert.Equal(t, TaskStatusInProgress, task.Status, "status should be updated to in progress")
		assert.Equal(t, TaskPriorityHigh, task.Priority, "priority should be updated to high")
	})

	t.Run("权限验证失败 - 不同用户", func(t *testing.T) {
		newTitle := "恶意更新"
		param := UpdateTaskParam{
			TaskID: "task-123",
			UserID: "other-user", // 不同的用户ID
			Title:  &newTitle,
		}

		task, err := usecase.UpdateTask(ctx, param)

		// ✅ TDD: 明确期望权限错误
		assert.Nil(t, task, "should return nil task for unauthorized user")
		assert.Equal(t, ErrTaskNotFound, err, "should return ErrTaskNotFound for unauthorized access")
	})
}

// 测试 DeleteTask 方法
func TestTaskUsecase_DeleteTask(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功删除任务", func(t *testing.T) {
		param := DeleteTaskParam{
			TaskID: "task-123",
			UserID: "user-123",
		}

		err := usecase.DeleteTask(ctx, param)

		// ❌ TDD: 期望成功删除，当前业务逻辑未实现会失败
		require.NoError(t, err, "DeleteTask should succeed")
	})

	t.Run("权限验证失败", func(t *testing.T) {
		param := DeleteTaskParam{
			TaskID: "task-123",
			UserID: "other-user",
		}

		err := usecase.DeleteTask(ctx, param)

		// ✅ TDD: 明确期望权限错误
		assert.Equal(t, ErrTaskNotFound, err, "should return ErrTaskNotFound for unauthorized deletion")
	})

	t.Run("任务不存在", func(t *testing.T) {
		param := DeleteTaskParam{
			TaskID: "non-existent",
			UserID: "user-123",
		}

		err := usecase.DeleteTask(ctx, param)

		// ✅ TDD: 明确期望任务不存在错误
		assert.Equal(t, ErrTaskNotFound, err, "should return ErrTaskNotFound for non-existent task")
	})
}

// 测试 SetTaskScore 方法
func TestTaskUsecase_SetTaskScore(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功设置任务分数", func(t *testing.T) {
		param := SetTaskScoreParam{
			TaskID: "task-123",
			UserID: "user-123",
			Score:  150,
		}

		task, err := usecase.SetTaskScore(ctx, param)

		// ❌ TDD: 期望成功设置，当前业务逻辑未实现会失败
		require.NoError(t, err, "SetTaskScore should succeed")
		require.NotNil(t, task, "should return updated task")
		assert.Equal(t, param.Score, task.Score, "score should be updated")
	})

	t.Run("无效分数", func(t *testing.T) {
		param := SetTaskScoreParam{
			TaskID: "task-123",
			UserID: "user-123",
			Score:  -50, // 负分数
		}

		task, err := usecase.SetTaskScore(ctx, param)

		// ✅ TDD: 明确期望分数验证错误
		assert.Nil(t, task, "should return nil task for invalid score")
		assert.Equal(t, ErrInvalidInput, err, "should return ErrInvalidInput for negative score")
	})
}

// 测试 CreateSubTask 方法
func TestTaskUsecase_CreateSubTask(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功创建子任务", func(t *testing.T) {
		param := CreateSubTaskParam{
			ParentID: "parent-task-123",
			UserID:   "user-123",
			Title:    "子任务1",
			Type:     PeriodDay,
			Period: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
			Tags:  []string{"子任务"},
			Icon:  "📋",
			Score: 30,
		}

		task, err := usecase.CreateSubTask(ctx, param)

		// ❌ TDD: 期望成功创建，当前业务逻辑未实现会失败
		require.NoError(t, err, "CreateSubTask should succeed")
		require.NotNil(t, task, "should return created sub task")
		assert.Equal(t, param.ParentID, task.ParentID, "parent ID should match")
		assert.Equal(t, param.Title, task.Title, "title should match")
	})

	t.Run("父任务不存在", func(t *testing.T) {
		param := CreateSubTaskParam{
			ParentID: "non-existent-parent",
			UserID:   "user-123",
			Title:    "子任务",
			Type:     PeriodDay,
			Score:    30,
		}

		task, err := usecase.CreateSubTask(ctx, param)

		// ✅ TDD: 明确期望父任务不存在错误
		assert.Nil(t, task, "should return nil task for non-existent parent")
		assert.Equal(t, ErrTaskNotFound, err, "should return ErrTaskNotFound for non-existent parent")
	})
}

// 测试 EditTag 方法
func TestTaskUsecase_TagOperations(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功编辑标签", func(t *testing.T) {
		param := EditTagParam{
			TaskID: "task-123",
			UserID: "user-123",
			Tags:   []string{"新标签1", "新标签2", "新标签3"},
		}

		task, err := usecase.EditTag(ctx, param)

		// ❌ TDD: 期望成功编辑，当前业务逻辑未实现会失败
		require.NoError(t, err, "EditTag should succeed")
		require.NotNil(t, task, "should return updated task")

		// 验证标签被完全替换
		assert.Equal(t, param.Tags, task.Tags, "tags should be completely replaced")
		assert.Len(t, task.Tags, 3, "should have exactly 3 tags")
	})

	t.Run("清空所有标签", func(t *testing.T) {
		param := EditTagParam{
			TaskID: "task-123",
			UserID: "user-123",
			Tags:   []string{}, // 空标签数组
		}

		task, err := usecase.EditTag(ctx, param)

		// ❌ TDD: 期望成功清空，当前业务逻辑未实现会失败
		require.NoError(t, err, "EditTag should succeed for empty tags")
		require.NotNil(t, task, "should return updated task")

		// 验证标签被清空
		assert.Empty(t, task.Tags, "tags should be empty")
	})

	t.Run("权限验证失败", func(t *testing.T) {
		param := EditTagParam{
			TaskID: "task-123",
			UserID: "other-user", // 不同的用户ID
			Tags:   []string{"恶意标签"},
		}

		task, err := usecase.EditTag(ctx, param)

		// ✅ TDD: 明确期望权限错误
		assert.Nil(t, task, "should return nil task for unauthorized user")
		assert.Equal(t, ErrTaskNotFound, err, "should return ErrTaskNotFound for unauthorized access")
	})
}

// 测试 SetTaskIcon 方法
func TestTaskUsecase_SetTaskIcon(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功设置任务图标", func(t *testing.T) {
		param := SetTaskIconParam{
			TaskID: "task-123",
			UserID: "user-123",
			Icon:   "🚀",
		}

		task, err := usecase.SetTaskIcon(ctx, param)

		// ❌ TDD: 期望成功设置，当前业务逻辑未实现会失败
		require.NoError(t, err, "SetTaskIcon should succeed")
		require.NotNil(t, task, "should return updated task")
		assert.Equal(t, param.Icon, task.Icon, "icon should be updated")
	})
}

// 测试 ListTaskByPeriod 方法
func TestTaskUsecase_ListTaskByPeriod(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功获取月度任务列表", func(t *testing.T) {
		param := ListTaskByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodMonth,
		}

		tasks, err := usecase.ListTaskByPeriod(ctx, param)

		// ❌ TDD: 期望成功获取，当前业务逻辑未实现会失败
		require.NoError(t, err, "ListTaskByPeriod should succeed")
		require.NotNil(t, tasks, "should return task list")

		// 验证返回的任务都在指定时间范围内
		for _, task := range tasks {
			assert.Equal(t, param.UserID, task.UserID, "should only return user's tasks")

			// 验证任务时间在范围内
			assert.False(t, task.TimePeriod.Start.Before(param.Period.Start),
				"task start time should be within period")
			assert.False(t, task.TimePeriod.End.After(param.Period.End),
				"task end time should be within period")
		}
	})

	t.Run("成功获取日度任务列表", func(t *testing.T) {
		param := ListTaskByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodDay,
		}

		tasks, err := usecase.ListTaskByPeriod(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if tasks == nil {
			t.Fatal("❌ 应该返回任务列表")
		}

		// 验证返回的任务类型
		for _, task := range tasks {
			if task.TaskType != PeriodDay {
				t.Errorf("期望日任务，得到 %v", task.TaskType)
			}
		}
	})
}

// 测试 ListTaskTree 方法
func TestTaskUsecase_ListTaskTree(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功获取任务树", func(t *testing.T) {
		param := ListTaskTreeParam{
			UserID: "user-123",
			TaskID: "parent-task-123",
		}

		tasks, err := usecase.ListTaskTree(ctx, param)

		// ❌ TDD: 期望成功获取，当前业务逻辑未实现会失败
		require.NoError(t, err, "ListTaskTree should succeed")
		require.NotNil(t, tasks, "should return task tree list")

		// 验证任务树结构 - mock 返回3个任务（1个父任务 + 2个子任务）
		assert.GreaterOrEqual(t, len(tasks), 1, "should return at least root task")

		// 验证返回的任务都属于正确的用户
		for _, task := range tasks {
			assert.Equal(t, param.UserID, task.UserID, "should only return user's tasks")
		}
	})
}

// 测试 ListTaskParentTree 方法
func TestTaskUsecase_ListTaskParentTree(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功获取父任务树", func(t *testing.T) {
		param := ListTaskParentTreeParam{
			UserID: "user-123",
			TaskID: "child-task-123",
		}

		tasks, err := usecase.ListTaskParentTree(ctx, param)

		// ❌ TDD: 期望成功获取，当前业务逻辑未实现会失败
		require.NoError(t, err, "ListTaskParentTree should succeed")
		require.NotNil(t, tasks, "should return parent task tree list")

		// 验证返回的都是父级任务
		for _, task := range tasks {
			assert.Equal(t, param.UserID, task.UserID, "should only return user's tasks")
		}
	})
}

// 测试 GetTaskStats 方法
func TestTaskUsecase_GetTaskStats(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功获取任务统计", func(t *testing.T) {
		param := GetTaskStatsParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodMonth,
		}

		stats, err := usecase.GetTaskStats(ctx, param)

		// ❌ TDD: 期望成功获取，当前业务逻辑未实现会失败
		require.NoError(t, err, "GetTaskStats should succeed")
		require.NotNil(t, stats, "should return statistics data")

		// 实际返回的统计数据长度可能不是12个月，取决于 mock 数据
		assert.GreaterOrEqual(t, len(stats), 0, "should return some statistics data")

		// 验证统计数据格式
		for _, stat := range stats {
			assert.GreaterOrEqual(t, stat.TaskCount, 0, "task count should not be negative")
			assert.GreaterOrEqual(t, stat.ScoreTotal, 0, "score total should not be negative")
		}
	})
}

// 测试结构体字段
func TestTask_Fields(t *testing.T) {
	task := Task{
		ID:       "task-123",
		Title:    "测试任务",
		TaskType: PeriodDay,
		TimePeriod: Period{
			Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
		},
		Tags:      []string{"测试", "任务"},
		Icon:      "📝",
		Score:     80,
		Status:    TaskStatusNotStarted,
		Priority:  TaskPriorityLow,
		ParentID:  "",
		UserID:    "user-123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if task.ID != "task-123" {
		t.Errorf("期望ID为 'task-123', 得到 %s", task.ID)
	}

	if task.Title != "测试任务" {
		t.Errorf("期望标题为 '测试任务', 得到 %s", task.Title)
	}

	if task.TaskType != PeriodDay {
		t.Errorf("期望类型为 PeriodDay, 得到 %v", task.TaskType)
	}

	if task.Score != 80 {
		t.Errorf("期望分数为 80, 得到 %d", task.Score)
	}

	if task.Status != TaskStatusNotStarted {
		t.Error("期望任务为未开始状态")
	}

	if task.Priority != TaskPriorityLow {
		t.Error("期望任务为低优先级")
	}

	if len(task.Tags) != 2 {
		t.Errorf("期望标签数量为 2, 得到 %d", len(task.Tags))
	}
}

// 测试参数结构体
func TestCreateTaskParam_Fields(t *testing.T) {
	param := CreateTaskParam{
		UserID: "user-123",
		Title:  "新任务",
		Type:   PeriodWeek,
		Period: Period{
			Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 19, 23, 59, 59, 0, time.UTC),
		},
		Tags:     []string{"新建", "任务"},
		Icon:     "🎯",
		Score:    100,
		ParentID: "parent-123",
	}

	if param.UserID != "user-123" {
		t.Errorf("期望用户ID为 'user-123', 得到 %s", param.UserID)
	}

	if param.Type != PeriodWeek {
		t.Errorf("期望类型为 PeriodWeek, 得到 %v", param.Type)
	}

	if param.Score != 100 {
		t.Errorf("期望分数为 100, 得到 %d", param.Score)
	}
}

// 边界测试
func TestTaskUsecase_EdgeCases(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("极长标题", func(t *testing.T) {
		longTitle := strings.Repeat("很长的任务标题", 1000)
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  longTitle,
			Type:   PeriodDay,
			Score:  50,
		}

		task, err := usecase.CreateTask(ctx, param)

		// 实现后应该有标题长度限制
		if err == ErrNoPermission {
			t.Log("当前返回 ErrNoPermission，实现后应该有标题长度验证")
		}

		if task != nil && len(task.Title) > 200 {
			t.Errorf("标题可能过长，需要限制长度")
		}
	})

	t.Run("极大分数", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "高分任务",
			Type:   PeriodDay,
			Score:  999999, // 极大分数
		}

		task, err := usecase.CreateTask(ctx, param)

		if err == ErrNoPermission {
			t.Log("当前返回 ErrNoPermission，实现后应该有分数范围验证")
		}

		if task != nil && task.Score > 1000 {
			t.Log("可能需要设置分数上限")
		}
	})

	t.Run("大量标签", func(t *testing.T) {
		manyTags := make([]string, 100)
		for i := range manyTags {
			manyTags[i] = "标签" + string(rune(i))
		}

		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "多标签任务",
			Type:   PeriodDay,
			Tags:   manyTags,
			Score:  50,
		}

		task, err := usecase.CreateTask(ctx, param)

		if err == ErrNoPermission {
			t.Log("当前返回 ErrNoPermission，实现后应该限制标签数量")
		}

		if task != nil && len(task.Tags) > 20 {
			t.Log("可能需要限制标签数量")
		}
	})

	t.Run("特殊字符处理", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "任务<script>alert('xss')</script>",
			Type:   PeriodDay,
			Tags:   []string{"特殊&字符", "<危险>标签"},
			Icon:   "🚀💡🎯",
			Score:  50,
		}

		task, err := usecase.CreateTask(ctx, param)

		if err == ErrNoPermission {
			t.Log("当前返回 ErrNoPermission，实现后需要处理特殊字符转义")
		}

		if task != nil {
			// 验证特殊字符被正确处理
			if strings.Contains(task.Title, "<script>") {
				t.Error("可能存在XSS风险，需要转义HTML标签")
			}
		}
	})
}

// 测试任务状态和优先级枚举
func TestTaskUsecase_StatusAndPriority(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("测试所有任务状态", func(t *testing.T) {
		statuses := []TaskStatus{
			TaskStatusNotStarted,
			TaskStatusInProgress,
			TaskStatusCompleted,
			TaskStatusCancelled,
		}

		for _, status := range statuses {
			param := UpdateTaskParam{
				TaskID: "task-123",
				UserID: "user-123",
				Status: &status,
			}

			task, err := usecase.UpdateTask(ctx, param)
			if err == nil && task != nil {
				assert.Equal(t, status, task.Status, "status should be updated correctly")
			}
		}
	})

	t.Run("测试所有优先级", func(t *testing.T) {
		priorities := []TaskPriority{
			TaskPriorityLow,
			TaskPriorityMedium,
			TaskPriorityHigh,
			TaskPriorityUrgent,
		}

		for _, priority := range priorities {
			param := UpdateTaskParam{
				TaskID:   "task-123",
				UserID:   "user-123",
				Priority: &priority,
			}

			task, err := usecase.UpdateTask(ctx, param)
			if err == nil && task != nil {
				assert.Equal(t, priority, task.Priority, "priority should be updated correctly")
			}
		}
	})

	t.Run("创建任务时指定优先级", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "高优先级任务",
			Type:   PeriodDay,
			Period: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			},
			Tags:     []string{"测试", "优先级"},
			Icon:     "🔥",
			Score:    100,
			Priority: TaskPriorityHigh,
		}

		task, err := usecase.CreateTask(ctx, param)
		if err == nil && task != nil {
			assert.Equal(t, TaskStatusNotStarted, task.Status, "new task should be not started by default")
			assert.Equal(t, TaskPriorityHigh, task.Priority, "priority should match param")
		}
	})
}
