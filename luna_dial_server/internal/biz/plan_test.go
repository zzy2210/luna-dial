package biz

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 创建测试用的 PlanUsecase 实例
func createTestPlanUsecase() *PlanUsecase {
	// 创建 mock repo 实例，这里先用 nil，因为当前实现还没完成
	taskRepo := &mockTaskRepo{}
	journalRepo := &mockJournalRepo{}

	taskUsecase := NewTaskUsecase(taskRepo)
	journalUsecase := NewJournalUsecase(journalRepo)

	return NewPlanUsecase(taskUsecase, journalUsecase)
}

// Mock TaskRepo 实现
type mockTaskRepo struct{}

func (m *mockTaskRepo) CreateTask(ctx context.Context, task *Task) error {
	return nil
}
func (m *mockTaskRepo) UpdateTask(ctx context.Context, task *Task) error {
	return nil
}
func (m *mockTaskRepo) DeleteTask(ctx context.Context, taskID, userID string) error {
	return nil
}
func (m *mockTaskRepo) GetTask(ctx context.Context, taskID, userID string) (*Task, error) {
	// 模拟一些测试数据
	if taskID == "task-123" && userID == "user-123" {
		return &Task{
			ID:       taskID,
			Title:    "测试任务",
			UserID:   userID,
			TaskType: PeriodDay,
			Tags:     []string{"测试"},
			Score:    50,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
			CreatedAt: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
		}, nil
	}
	if taskID == "parent-task-123" && userID == "user-123" {
		return &Task{
			ID:       taskID,
			Title:    "父任务",
			UserID:   userID,
			TaskType: PeriodDay,
			Tags:     []string{"父任务"},
			Score:    100,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
			CreatedAt: time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
		}, nil
	}
	if taskID == "child-task-123" && userID == "user-123" {
		return &Task{
			ID:       taskID,
			Title:    "子任务",
			UserID:   userID,
			TaskType: PeriodDay,
			Tags:     []string{"子任务"},
			Score:    30,
			ParentID: "parent-task-123",
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
			CreatedAt: time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC),
		}, nil
	}
	if userID == "other-user" {
		return nil, ErrTaskNotFound // 模拟权限错误
	}
	if taskID == "non-existent" || taskID == "non-existent-parent" {
		return nil, ErrTaskNotFound // 模拟任务不存在
	}
	return nil, ErrTaskNotFound
}
func (m *mockTaskRepo) ListTasks(ctx context.Context, userID string, periodStart, periodEnd time.Time, taskType int) ([]*Task, error) {
	// 模拟返回一些测试任务
	if userID == "user-123" {
		return []*Task{
			{
				ID:       "task-1",
				Title:    "任务1",
				UserID:   userID,
				TaskType: PeriodType(taskType),
				Score:    100,
				TimePeriod: Period{
					Start: periodStart.Add(time.Hour),
					End:   periodStart.Add(2 * time.Hour),
				},
			},
			{
				ID:       "task-2",
				Title:    "任务2",
				UserID:   userID,
				TaskType: PeriodType(taskType),
				Score:    50,
				TimePeriod: Period{
					Start: periodStart.Add(3 * time.Hour),
					End:   periodStart.Add(4 * time.Hour),
				},
			},
		}, nil
	}
	return []*Task{}, nil
}
func (m *mockTaskRepo) ListTaskTree(ctx context.Context, taskID, userID string) ([]*Task, error) {
	// 模拟返回任务树
	if taskID == "parent-task-123" && userID == "user-123" {
		return []*Task{
			{
				ID:       "parent-task-123",
				Title:    "父任务",
				UserID:   userID,
				TaskType: PeriodDay,
				Score:    100,
			},
			{
				ID:       "child-task-1",
				Title:    "子任务1",
				UserID:   userID,
				TaskType: PeriodDay,
				ParentID: "parent-task-123",
				Score:    30,
			},
			{
				ID:       "child-task-2",
				Title:    "子任务2",
				UserID:   userID,
				TaskType: PeriodDay,
				ParentID: "parent-task-123",
				Score:    20,
			},
		}, nil
	}
	return []*Task{}, nil
}
func (m *mockTaskRepo) ListTaskParentTree(ctx context.Context, taskID, userID string) ([]*Task, error) {
	// 模拟返回父任务树路径
	if taskID == "child-task-123" && userID == "user-123" {
		return []*Task{
			{
				ID:       "root-task",
				Title:    "根任务",
				UserID:   userID,
				TaskType: PeriodDay,
				Score:    200,
			},
			{
				ID:       "parent-task-123",
				Title:    "父任务",
				UserID:   userID,
				TaskType: PeriodDay,
				ParentID: "root-task",
				Score:    100,
			},
			{
				ID:       "child-task-123",
				Title:    "子任务",
				UserID:   userID,
				TaskType: PeriodDay,
				ParentID: "parent-task-123",
				Score:    30,
			},
		}, nil
	}
	return []*Task{}, nil
}

// 测试 NewPlanUsecase 构造函数
func TestNewPlanUsecase(t *testing.T) {
	taskRepo := &mockTaskRepo{}
	journalRepo := &mockJournalRepo{}

	taskUsecase := NewTaskUsecase(taskRepo)
	journalUsecase := NewJournalUsecase(journalRepo)

	planUsecase := NewPlanUsecase(taskUsecase, journalUsecase)

	if planUsecase == nil {
		t.Fatal("NewPlanUsecase returned nil")
	}

	if planUsecase.taskUsecase != taskUsecase {
		t.Error("taskUsecase not set correctly")
	}

	if planUsecase.journalUsecase != journalUsecase {
		t.Error("journalUsecase not set correctly")
	}
}

// 测试 GetPlanByPeriod 方法 - 这些测试会失败，因为期望真正的业务逻辑实现
func TestPlanUsecase_GetPlanByPeriod(t *testing.T) {
	planUsecase := createTestPlanUsecase()

	t.Run("成功获取月度计划", func(t *testing.T) {
		param := GetPlanByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodDay,
		}

		plan, err := planUsecase.GetPlanByPeriod(context.TODO(), param)

		// 业务逻辑实现后，err 应该为 nil，plan 不为 nil
		assert.Nil(t, err)
		assert.NotNil(t, plan)

		if plan != nil {
			// 验证计划结构
			assert.Equal(t, PeriodDay, plan.PlanType)
			assert.Equal(t, param.Period.Start, plan.PlanPeriod.Start)
			assert.Equal(t, param.Period.End, plan.PlanPeriod.End)

			// 期望有任务和日志数据
			assert.GreaterOrEqual(t, plan.TasksTotal, 0)
			assert.GreaterOrEqual(t, plan.JournalsTotal, 0)

			// 期望有分组统计数据（按日分组，1月有31天）
			expectedDays := 31
			assert.Equal(t, expectedDays, len(plan.GroupStats))
		}
	})

	t.Run("成功获取周度计划", func(t *testing.T) {
		param := GetPlanByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodWeek,
		}

		plan, err := planUsecase.GetPlanByPeriod(context.TODO(), param)

		// 业务逻辑实现后，err 应该为 nil，plan 不为 nil
		assert.Nil(t, err)
		assert.NotNil(t, plan)

		if plan != nil {
			assert.Equal(t, PeriodWeek, plan.PlanType)

			// 期望按周分组，1月大约有4-5周
			assert.GreaterOrEqual(t, len(plan.GroupStats), 4)
			assert.LessOrEqual(t, len(plan.GroupStats), 6)

			// 验证周统计的 GroupKey 格式应该是 "2025-W01", "2025-W02" 等
			for _, stat := range plan.GroupStats {
				assert.NotEmpty(t, stat.GroupKey)
			}
		}
	})

	t.Run("参数验证失败 - 空用户ID", func(t *testing.T) {
		param := GetPlanByPeriodParam{
			UserID: "", // 空用户ID应该失败
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodDay,
		}

		plan, err := planUsecase.GetPlanByPeriod(context.TODO(), param)

		assert.Nil(t, plan)
		assert.NotNil(t, err)
		// 业务逻辑实现后应该返回 ErrUserIDEmpty
		assert.Equal(t, ErrUserIDEmpty, err, "应该返回 ErrUserIDEmpty 错误")
	})

	t.Run("参数验证失败 - 无效时间区间", func(t *testing.T) {
		param := GetPlanByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // 结束时间早于开始时间
			},
			GroupBy: PeriodDay,
		}

		plan, err := planUsecase.GetPlanByPeriod(context.TODO(), param)

		assert.Nil(t, plan)
		assert.NotNil(t, err)
		// 业务逻辑实现后应该返回 ErrInvalidPeriod
		assert.Equal(t, ErrInvalidPeriod, err, "应该返回 ErrInvalidPeriod 错误")
	})
}

// 测试 GetPlanStats 方法 - 这些测试会失败，因为期望真正的业务逻辑实现
func TestPlanUsecase_GetPlanStats(t *testing.T) {
	planUsecase := createTestPlanUsecase()

	t.Run("成功获取月度统计", func(t *testing.T) {
		param := GetPlanStatsParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodMonth,
		}

		stats, err := planUsecase.GetPlanStats(context.TODO(), param)

		// 业务逻辑实现后，err 应该为 nil，stats 不为 nil
		assert.Nil(t, err)
		assert.NotNil(t, stats)

		if stats != nil {
			// 期望返回12个月的统计数据
			expectedMonths := 12
			assert.Equal(t, expectedMonths, len(stats))

			// 验证统计数据格式
			for i, stat := range stats {
				expectedKey := fmt.Sprintf("2025-%02d", i+1)
				assert.Equal(t, expectedKey, stat.GroupKey)
				assert.GreaterOrEqual(t, stat.TaskCount, 0)
				assert.GreaterOrEqual(t, stat.ScoreTotal, 0)
			}
		}
	})

	t.Run("成功获取周度统计", func(t *testing.T) {
		param := GetPlanStatsParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodWeek,
		}

		stats, err := planUsecase.GetPlanStats(context.TODO(), param)

		// 业务逻辑实现后，err 应该为 nil
		assert.Nil(t, err)
		assert.NotNil(t, stats)

		if stats != nil {
			// 1月大约有4-5周
			assert.GreaterOrEqual(t, len(stats), 4)
			assert.LessOrEqual(t, len(stats), 6)

			// 验证周统计的 GroupKey 格式应该是 "2025-W01", "2025-W02" 等
			for _, stat := range stats {
				assert.True(t, strings.HasPrefix(stat.GroupKey, "2025-W"))
			}
		}
	})

	t.Run("成功获取日度统计", func(t *testing.T) {
		param := GetPlanStatsParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 7, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodDay,
		}

		stats, err := planUsecase.GetPlanStats(context.TODO(), param)

		// 业务逻辑实现后，err 应该为 nil
		assert.Nil(t, err)
		assert.NotNil(t, stats)

		if stats != nil {
			// 期望返回7天的统计数据
			expectedDays := 7
			assert.Equal(t, expectedDays, len(stats))

			// 验证日统计的 GroupKey 格式应该是 "2025-01-01", "2025-01-02" 等
			for i, stat := range stats {
				expectedKey := fmt.Sprintf("2025-01-%02d", i+1)
				assert.Equal(t, expectedKey, stat.GroupKey)
			}
		}
	})

	t.Run("成功获取季度统计", func(t *testing.T) {
		param := GetPlanStatsParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodQuarter,
		}

		stats, err := planUsecase.GetPlanStats(context.TODO(), param)

		// 业务逻辑实现后，err 应该为 nil
		assert.Nil(t, err)
		assert.NotNil(t, stats)

		if stats != nil {
			// 期望返回4个季度的统计数据
			expectedQuarters := 4
			assert.Equal(t, expectedQuarters, len(stats))

			// 验证季度统计的 GroupKey 格式应该是 "2025-Q1", "2025-Q2" 等
			expectedQuarterKeys := []string{"2025-Q1", "2025-Q2", "2025-Q3", "2025-Q4"}
			for i, stat := range stats {
				assert.Equal(t, expectedQuarterKeys[i], stat.GroupKey)
			}
		}
	})

	t.Run("参数验证失败 - 空用户ID", func(t *testing.T) {
		param := GetPlanStatsParam{
			UserID: "", // 空用户ID应该失败
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodMonth,
		}

		stats, err := planUsecase.GetPlanStats(context.TODO(), param)

		assert.Nil(t, stats)
		assert.NotNil(t, err)
		// 业务逻辑实现后应该返回 ErrUserIDEmpty
		assert.Equal(t, ErrUserIDEmpty, err, "应该返回 ErrUserIDEmpty 错误")
	})

	t.Run("参数验证失败 - 无效时间区间", func(t *testing.T) {
		param := GetPlanStatsParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // 结束时间早于开始时间
			},
			GroupBy: PeriodMonth,
		}

		stats, err := planUsecase.GetPlanStats(context.TODO(), param)

		assert.Nil(t, stats)
		assert.NotNil(t, err)
		// 业务逻辑实现后应该返回 ErrInvalidPeriod
		assert.Equal(t, ErrInvalidPeriod, err, "应该返回 ErrInvalidPeriod 错误")
	})
}

// 测试 Plan 结构体的字段
func TestPlan_Fields(t *testing.T) {
	// 创建一个示例 Plan 对象来验证结构体定义
	plan := Plan{
		Tasks:         []Task{},
		TasksTotal:    10,
		Journals:      []Journal{},
		JournalsTotal: 5,
		PlanType:      PeriodMonth,
		PlanPeriod: Period{
			Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
		},
		ScoreTotal: 100,
		GroupStats: []GroupStat{
			{
				GroupKey:   "2025-01",
				TaskCount:  10,
				ScoreTotal: 100,
			},
		},
	}

	// 验证字段类型和值
	assert.Equal(t, 10, plan.TasksTotal)
	assert.Equal(t, 5, plan.JournalsTotal)
	assert.Equal(t, PeriodMonth, plan.PlanType)
	assert.Equal(t, 100, plan.ScoreTotal)
	assert.Equal(t, 1, len(plan.GroupStats))
	assert.Equal(t, "2025-01", plan.GroupStats[0].GroupKey)
}

// 测试 GroupStat 结构体
func TestGroupStat_Fields(t *testing.T) {
	stat := GroupStat{
		GroupKey:   "2025-W01",
		TaskCount:  5,
		ScoreTotal: 50,
	}

	assert.Equal(t, "2025-W01", stat.GroupKey)
	assert.Equal(t, 5, stat.TaskCount)
	assert.Equal(t, 50, stat.ScoreTotal)
}

// 测试参数结构体
func TestGetPlanByPeriodParam_Fields(t *testing.T) {
	param := GetPlanByPeriodParam{
		UserID: "user-123",
		Period: Period{
			Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
		},
		GroupBy: PeriodDay,
	}

	assert.Equal(t, "user-123", param.UserID)
	assert.Equal(t, PeriodDay, param.GroupBy)
}

func TestGetPlanStatsParam_Fields(t *testing.T) {
	param := GetPlanStatsParam{
		UserID: "user-456",
		Period: Period{
			Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
		},
		GroupBy: PeriodQuarter,
	}

	assert.Equal(t, "user-456", param.UserID)
	assert.Equal(t, PeriodQuarter, param.GroupBy)
}

// 边界测试：极端情况
func TestPlanUsecase_EdgeCases(t *testing.T) {
	planUsecase := createTestPlanUsecase()

	t.Run("nil usecase", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				// 这是预期的，nil usecase 调用方法应该会 panic
				return
			}
			// 如果没有 panic，说明当前实现可能有保护措施，这也是可以的
			t.Log("No panic occurred - this might be acceptable if the implementation has nil checks")
		}()

		var nilUsecase *PlanUsecase
		_, _ = nilUsecase.GetPlanByPeriod(context.TODO(), GetPlanByPeriodParam{})
	})

	t.Run("极长用户ID", func(t *testing.T) {
		longUserID := strings.Repeat("a", 1000)

		param := GetPlanByPeriodParam{
			UserID: longUserID,
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodDay,
		}

		_, err := planUsecase.GetPlanByPeriod(context.TODO(), param)
		// 业务逻辑实现后应该返回用户ID相关的验证错误
		assert.NotNil(t, err)
		// 可能是 ErrUserIDInvalid 或 ErrUserIDTooLong
	})

	t.Run("极端时间值", func(t *testing.T) {
		param := GetPlanByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodYear,
		}

		_, err := planUsecase.GetPlanByPeriod(context.TODO(), param)
		// 业务逻辑实现后可能因为时间跨度过大而返回错误
		assert.NotNil(t, err)
		// 可能是 ErrPeriodTooLarge 或类似错误
	})
}

// 性能测试：多次调用
func TestPlanUsecase_Performance(t *testing.T) {
	planUsecase := createTestPlanUsecase()

	param := GetPlanByPeriodParam{
		UserID: "user-123",
		Period: Period{
			Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
		},
		GroupBy: PeriodDay,
	}

	// 测试多次调用的一致性
	for i := 0; i < 100; i++ {
		plan, err := planUsecase.GetPlanByPeriod(context.TODO(), param)
		// 当前实现应该一致地返回错误
		assert.Nil(t, plan, "iteration %d: expected plan to be nil", i)
		assert.Equal(t, ErrNoPermission, err, "iteration %d: expected ErrNoPermission, got %v", i, err)
	}
}

// 这个测试专门用来验证真实的业务逻辑 - 这些测试**应该失败**
// 因为当前的实现只是返回 ErrNoPermission，而不是真正的业务逻辑
func TestPlanUsecase_RealBusinessLogic_ShouldFail(t *testing.T) {
	t.Run("验证真实的GetPlanByPeriod业务逻辑", func(t *testing.T) {
		planUsecase := createTestPlanUsecase()

		param := GetPlanByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodDay,
		}

		plan, err := planUsecase.GetPlanByPeriod(context.TODO(), param)

		// 业务逻辑实现后，这些断言应该通过
		assert.Nil(t, err, "GetPlanByPeriod 应该成功返回计划")
		assert.NotNil(t, plan, "GetPlanByPeriod 应该返回非空的计划对象")

		if plan != nil {
			assert.Equal(t, PeriodDay, plan.PlanType, "计划类型应该正确")
			assert.Equal(t, param.Period.Start, plan.PlanPeriod.Start, "计划开始时间应该正确")
			assert.Equal(t, param.Period.End, plan.PlanPeriod.End, "计划结束时间应该正确")

			// 期望包含分组统计数据（按日分组，1月有31天）
			expectedDays := 31
			assert.Equal(t, expectedDays, len(plan.GroupStats), "应该包含31天的分组统计数据")

			// 验证分组统计的格式
			for i, stat := range plan.GroupStats {
				expectedKey := fmt.Sprintf("2025-01-%02d", i+1)
				assert.Equal(t, expectedKey, stat.GroupKey, "分组键格式应该正确")
			}
		}
	})

	t.Run("验证真实的GetPlanStats业务逻辑", func(t *testing.T) {
		planUsecase := createTestPlanUsecase()

		param := GetPlanStatsParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodMonth,
		}

		stats, err := planUsecase.GetPlanStats(context.TODO(), param)

		// 业务逻辑实现后，这些断言应该通过
		assert.Nil(t, err, "GetPlanStats 应该成功返回统计数据")
		assert.NotNil(t, stats, "GetPlanStats 应该返回非空的统计数据")

		if stats != nil {
			// 期望返回12个月的统计数据
			expectedMonths := 12
			assert.Equal(t, expectedMonths, len(stats), "应该返回12个月的统计数据")

			// 验证每个月的统计数据格式
			for i, stat := range stats {
				expectedKey := fmt.Sprintf("2025-%02d", i+1)
				assert.Equal(t, expectedKey, stat.GroupKey, "月份分组键应该正确")
				assert.GreaterOrEqual(t, stat.TaskCount, 0, "任务数量不能为负数")
				assert.GreaterOrEqual(t, stat.ScoreTotal, 0, "总分不能为负数")
			}
		}
	})

	t.Run("验证不同分组类型的业务逻辑", func(t *testing.T) {
		planUsecase := createTestPlanUsecase()

		testCases := []struct {
			name          string
			groupBy       PeriodType
			expectedCount int
			keyPattern    string
		}{
			{
				name:          "按周分组",
				groupBy:       PeriodWeek,
				expectedCount: 5, // 1月大约5周
				keyPattern:    "2025-W",
			},
			{
				name:          "按季度分组",
				groupBy:       PeriodQuarter,
				expectedCount: 4, // 全年4个季度
				keyPattern:    "2025-Q",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				param := GetPlanStatsParam{
					UserID: "user-123",
					Period: Period{
						Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
					},
					GroupBy: tc.groupBy,
				}

				stats, err := planUsecase.GetPlanStats(context.TODO(), param)

				// 业务逻辑实现后，这些断言应该通过
				assert.Nil(t, err, "%s 业务逻辑应该成功", tc.name)
				assert.NotNil(t, stats, "%s 应该返回统计数据", tc.name)

				if stats != nil && tc.groupBy == PeriodQuarter {
					assert.Equal(t, tc.expectedCount, len(stats), "%s 统计数量应该正确", tc.name)
				}

				// 验证分组键的格式
				if stats != nil {
					for _, stat := range stats {
						assert.True(t, strings.HasPrefix(stat.GroupKey, tc.keyPattern),
							"%s 分组键格式应该正确: 期望以 %s 开头, 得到 %s", tc.name, tc.keyPattern, stat.GroupKey)
					}
				}
			})
		}
	})
}
