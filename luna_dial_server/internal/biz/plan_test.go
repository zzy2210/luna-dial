package biz

import (
	"fmt"
	"strings"
	"testing"
	"time"
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

func (m *mockTaskRepo) CreateTask(task *Task) error                  { return nil }
func (m *mockTaskRepo) UpdateTask(task *Task) error                  { return nil }
func (m *mockTaskRepo) DeleteTask(taskID, userID string) error       { return nil }
func (m *mockTaskRepo) GetTask(taskID, userID string) (*Task, error) { return nil, ErrTaskNotFound }
func (m *mockTaskRepo) ListTasks(userID string, periodStart, periodEnd time.Time, taskType string) ([]*Task, error) {
	return nil, nil
}
func (m *mockTaskRepo) ListTaskTree(taskID, userID string) ([]*Task, error) { return nil, nil }

// Mock JournalRepo 实现
type mockJournalRepo struct{}

func (m *mockJournalRepo) CreateJournal(journal *Journal) error         { return nil }
func (m *mockJournalRepo) UpdateJournal(journal *Journal) error         { return nil }
func (m *mockJournalRepo) DeleteJournal(journalID, userID string) error { return nil }
func (m *mockJournalRepo) GetJournal(journalID, userID string) (*Journal, error) {
	return nil, ErrJournalNotFound
}
func (m *mockJournalRepo) ListJournals(userID string, periodStart, periodEnd time.Time, journalType string) ([]*Journal, error) {
	return nil, nil
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

		plan, err := planUsecase.GetPlanByPeriod(param)

		// 期望成功返回计划，但实际会失败因为当前返回 ErrNoPermission
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if plan == nil {
			t.Fatal("expected plan to be returned, got nil")
		}

		// 验证计划结构
		if plan.PlanType != PeriodDay {
			t.Errorf("expected PlanType to be PeriodDay, got %v", plan.PlanType)
		}

		if plan.PlanPeriod.Start != param.Period.Start {
			t.Errorf("expected plan start time to match param")
		}

		if plan.PlanPeriod.End != param.Period.End {
			t.Errorf("expected plan end time to match param")
		}

		// 期望有任务和日志数据
		if plan.TasksTotal < 0 {
			t.Errorf("expected non-negative TasksTotal")
		}

		if plan.JournalsTotal < 0 {
			t.Errorf("expected non-negative JournalsTotal")
		}

		// 期望有分组统计数据
		expectedDays := 31 // 1月有31天
		if len(plan.GroupStats) != expectedDays {
			t.Errorf("expected %d group stats for daily grouping, got %d", expectedDays, len(plan.GroupStats))
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

		plan, err := planUsecase.GetPlanByPeriod(param)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if plan == nil {
			t.Fatal("expected plan to be returned, got nil")
		}

		if plan.PlanType != PeriodWeek {
			t.Errorf("expected PlanType to be PeriodWeek, got %v", plan.PlanType)
		}

		// 期望按周分组，1月大约有4-5周
		if len(plan.GroupStats) < 4 || len(plan.GroupStats) > 6 {
			t.Errorf("expected 4-6 group stats for weekly grouping, got %d", len(plan.GroupStats))
		}

		// 验证周统计的 GroupKey 格式应该是 "2025-W01", "2025-W02" 等
		for _, stat := range plan.GroupStats {
			if len(stat.GroupKey) == 0 {
				t.Errorf("expected non-empty GroupKey")
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

		plan, err := planUsecase.GetPlanByPeriod(param)

		if plan != nil {
			t.Errorf("expected plan to be nil for invalid input, got %+v", plan)
		}

		// 期望返回特定的验证错误，而不是 ErrNoPermission
		if err == nil {
			t.Error("expected error for empty user ID, got nil")
		}

		// TODO: 实现后应该返回 ErrUserIDEmpty 或类似错误
		if err == ErrNoPermission {
			t.Log("当前实现返回 ErrNoPermission，实现后应该返回更具体的验证错误")
		}
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

		plan, err := planUsecase.GetPlanByPeriod(param)

		if plan != nil {
			t.Errorf("expected plan to be nil for invalid period, got %+v", plan)
		}

		if err == nil {
			t.Error("expected error for invalid period, got nil")
		}

		// TODO: 实现后应该返回 ErrInvalidPeriod 或类似错误
		if err == ErrNoPermission {
			t.Log("当前实现返回 ErrNoPermission，实现后应该返回更具体的验证错误")
		}
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

		stats, err := planUsecase.GetPlanStats(param)

		// 期望成功返回统计数据，但实际会失败因为当前返回 ErrNoPermission
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if stats == nil {
			t.Fatal("expected stats to be returned, got nil")
		}

		// 期望返回12个月的统计数据
		expectedMonths := 12
		if len(stats) != expectedMonths {
			t.Errorf("expected %d month stats, got %d", expectedMonths, len(stats))
		}

		// 验证统计数据格式
		for i, stat := range stats {
			expectedKey := fmt.Sprintf("2025-%02d", i+1)
			if stat.GroupKey != expectedKey {
				t.Errorf("expected GroupKey to be %s, got %s", expectedKey, stat.GroupKey)
			}

			if stat.TaskCount < 0 {
				t.Errorf("expected non-negative TaskCount, got %d", stat.TaskCount)
			}

			if stat.ScoreTotal < 0 {
				t.Errorf("expected non-negative ScoreTotal, got %d", stat.ScoreTotal)
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

		stats, err := planUsecase.GetPlanStats(param)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if stats == nil {
			t.Fatal("expected stats to be returned, got nil")
		}

		// 1月大约有4-5周
		if len(stats) < 4 || len(stats) > 6 {
			t.Errorf("expected 4-6 week stats, got %d", len(stats))
		}

		// 验证周统计的 GroupKey 格式应该是 "2025-W01", "2025-W02" 等
		for _, stat := range stats {
			if !strings.HasPrefix(stat.GroupKey, "2025-W") {
				t.Errorf("expected GroupKey to start with '2025-W', got %s", stat.GroupKey)
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

		stats, err := planUsecase.GetPlanStats(param)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if stats == nil {
			t.Fatal("expected stats to be returned, got nil")
		}

		// 期望返回7天的统计数据
		expectedDays := 7
		if len(stats) != expectedDays {
			t.Errorf("expected %d day stats, got %d", expectedDays, len(stats))
		}

		// 验证日统计的 GroupKey 格式应该是 "2025-01-01", "2025-01-02" 等
		for i, stat := range stats {
			expectedKey := fmt.Sprintf("2025-01-%02d", i+1)
			if stat.GroupKey != expectedKey {
				t.Errorf("expected GroupKey to be %s, got %s", expectedKey, stat.GroupKey)
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

		stats, err := planUsecase.GetPlanStats(param)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if stats == nil {
			t.Fatal("expected stats to be returned, got nil")
		}

		// 期望返回4个季度的统计数据
		expectedQuarters := 4
		if len(stats) != expectedQuarters {
			t.Errorf("expected %d quarter stats, got %d", expectedQuarters, len(stats))
		}

		// 验证季度统计的 GroupKey 格式应该是 "2025-Q1", "2025-Q2" 等
		expectedQuarterKeys := []string{"2025-Q1", "2025-Q2", "2025-Q3", "2025-Q4"}
		for i, stat := range stats {
			if stat.GroupKey != expectedQuarterKeys[i] {
				t.Errorf("expected GroupKey to be %s, got %s", expectedQuarterKeys[i], stat.GroupKey)
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

		stats, err := planUsecase.GetPlanStats(param)

		if stats != nil {
			t.Errorf("expected stats to be nil for invalid input, got %+v", stats)
		}

		if err == nil {
			t.Error("expected error for empty user ID, got nil")
		}

		// TODO: 实现后应该返回 ErrUserIDEmpty 或类似错误
		if err == ErrNoPermission {
			t.Log("当前实现返回 ErrNoPermission，实现后应该返回更具体的验证错误")
		}
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

		stats, err := planUsecase.GetPlanStats(param)

		if stats != nil {
			t.Errorf("expected stats to be nil for invalid period, got %+v", stats)
		}

		if err == nil {
			t.Error("expected error for invalid period, got nil")
		}

		// TODO: 实现后应该返回 ErrInvalidPeriod 或类似错误
		if err == ErrNoPermission {
			t.Log("当前实现返回 ErrNoPermission，实现后应该返回更具体的验证错误")
		}
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
	if plan.TasksTotal != 10 {
		t.Errorf("expected TasksTotal to be 10, got %d", plan.TasksTotal)
	}

	if plan.JournalsTotal != 5 {
		t.Errorf("expected JournalsTotal to be 5, got %d", plan.JournalsTotal)
	}

	if plan.PlanType != PeriodMonth {
		t.Errorf("expected PlanType to be PeriodMonth, got %v", plan.PlanType)
	}

	if plan.ScoreTotal != 100 {
		t.Errorf("expected ScoreTotal to be 100, got %d", plan.ScoreTotal)
	}

	if len(plan.GroupStats) != 1 {
		t.Errorf("expected GroupStats length to be 1, got %d", len(plan.GroupStats))
	}

	if plan.GroupStats[0].GroupKey != "2025-01" {
		t.Errorf("expected GroupKey to be '2025-01', got %s", plan.GroupStats[0].GroupKey)
	}
}

// 测试 GroupStat 结构体
func TestGroupStat_Fields(t *testing.T) {
	stat := GroupStat{
		GroupKey:   "2025-W01",
		TaskCount:  5,
		ScoreTotal: 50,
	}

	if stat.GroupKey != "2025-W01" {
		t.Errorf("expected GroupKey to be '2025-W01', got %s", stat.GroupKey)
	}

	if stat.TaskCount != 5 {
		t.Errorf("expected TaskCount to be 5, got %d", stat.TaskCount)
	}

	if stat.ScoreTotal != 50 {
		t.Errorf("expected ScoreTotal to be 50, got %d", stat.ScoreTotal)
	}
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

	if param.UserID != "user-123" {
		t.Errorf("expected UserID to be 'user-123', got %s", param.UserID)
	}

	if param.GroupBy != PeriodDay {
		t.Errorf("expected GroupBy to be PeriodDay, got %v", param.GroupBy)
	}
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

	if param.UserID != "user-456" {
		t.Errorf("expected UserID to be 'user-456', got %s", param.UserID)
	}

	if param.GroupBy != PeriodQuarter {
		t.Errorf("expected GroupBy to be PeriodQuarter, got %v", param.GroupBy)
	}
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
		_, _ = nilUsecase.GetPlanByPeriod(GetPlanByPeriodParam{})
	})

	t.Run("极长用户ID", func(t *testing.T) {
		longUserID := string(make([]byte, 1000))
		for i := range longUserID {
			longUserID = longUserID[:i] + "a" + longUserID[i+1:]
		}

		param := GetPlanByPeriodParam{
			UserID: longUserID,
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodDay,
		}

		_, err := planUsecase.GetPlanByPeriod(param)
		if err != ErrNoPermission {
			t.Errorf("expected ErrNoPermission, got %v", err)
		}
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

		_, err := planUsecase.GetPlanByPeriod(param)
		if err != ErrNoPermission {
			t.Errorf("expected ErrNoPermission, got %v", err)
		}
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
		plan, err := planUsecase.GetPlanByPeriod(param)
		if plan != nil {
			t.Errorf("iteration %d: expected plan to be nil", i)
		}
		if err != ErrNoPermission {
			t.Errorf("iteration %d: expected ErrNoPermission, got %v", i, err)
		}
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

		plan, err := planUsecase.GetPlanByPeriod(param)

		// 这个断言会失败，因为当前实现返回 nil 和 ErrNoPermission
		// 但我们期望的是真正的 Plan 对象
		if err != nil {
			t.Fatalf("❌ 业务逻辑未实现: GetPlanByPeriod 应该成功返回计划，但得到错误: %v", err)
		}

		if plan == nil {
			t.Fatal("❌ 业务逻辑未实现: GetPlanByPeriod 应该返回非空的计划对象")
		}

		// 验证返回的计划包含正确的字段
		if plan.PlanType != PeriodDay {
			t.Errorf("❌ 计划类型不正确: 期望 %v, 得到 %v", PeriodDay, plan.PlanType)
		}

		if plan.PlanPeriod.Start != param.Period.Start {
			t.Errorf("❌ 计划开始时间不正确")
		}

		if plan.PlanPeriod.End != param.Period.End {
			t.Errorf("❌ 计划结束时间不正确")
		}

		// 期望包含任务和日志数据
		if len(plan.Tasks) == 0 && plan.TasksTotal == 0 {
			t.Log("⚠️  当前没有任务数据，实现后应该从 TaskUsecase 获取")
		}

		if len(plan.Journals) == 0 && plan.JournalsTotal == 0 {
			t.Log("⚠️  当前没有日志数据，实现后应该从 JournalUsecase 获取")
		}

		// 期望包含分组统计数据（按日分组，1月有31天）
		expectedDays := 31
		if len(plan.GroupStats) != expectedDays {
			t.Errorf("❌ 分组统计数据不正确: 期望 %d 天的统计，得到 %d", expectedDays, len(plan.GroupStats))
		}

		// 验证分组统计的格式
		for i, stat := range plan.GroupStats {
			expectedKey := fmt.Sprintf("2025-01-%02d", i+1)
			if stat.GroupKey != expectedKey {
				t.Errorf("❌ 分组键格式不正确: 期望 %s, 得到 %s", expectedKey, stat.GroupKey)
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

		stats, err := planUsecase.GetPlanStats(param)

		// 这个断言会失败，因为当前实现返回 nil 和 ErrNoPermission
		if err != nil {
			t.Fatalf("❌ 业务逻辑未实现: GetPlanStats 应该成功返回统计数据，但得到错误: %v", err)
		}

		if stats == nil {
			t.Fatal("❌ 业务逻辑未实现: GetPlanStats 应该返回非空的统计数据")
		}

		// 期望返回12个月的统计数据
		expectedMonths := 12
		if len(stats) != expectedMonths {
			t.Errorf("❌ 统计数据数量不正确: 期望 %d 个月的统计，得到 %d", expectedMonths, len(stats))
		}

		// 验证每个月的统计数据格式
		for i, stat := range stats {
			expectedKey := fmt.Sprintf("2025-%02d", i+1)
			if stat.GroupKey != expectedKey {
				t.Errorf("❌ 月份分组键不正确: 期望 %s, 得到 %s", expectedKey, stat.GroupKey)
			}

			// 期望有有效的统计数据（任务数量和分数）
			if stat.TaskCount < 0 {
				t.Errorf("❌ 任务数量不能为负数: %d", stat.TaskCount)
			}

			if stat.ScoreTotal < 0 {
				t.Errorf("❌ 总分不能为负数: %d", stat.ScoreTotal)
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

				stats, err := planUsecase.GetPlanStats(param)

				if err != nil {
					t.Fatalf("❌ %s 业务逻辑未实现: %v", tc.name, err)
				}

				if stats == nil {
					t.Fatalf("❌ %s 应该返回统计数据", tc.name)
				}

				if tc.groupBy == PeriodQuarter && len(stats) != tc.expectedCount {
					t.Errorf("❌ %s 统计数量不正确: 期望 %d, 得到 %d", tc.name, tc.expectedCount, len(stats))
				}

				// 验证分组键的格式
				for _, stat := range stats {
					if !strings.HasPrefix(stat.GroupKey, tc.keyPattern) {
						t.Errorf("❌ %s 分组键格式不正确: 期望以 %s 开头, 得到 %s", tc.name, tc.keyPattern, stat.GroupKey)
					}
				}
			})
		}
	})
}
