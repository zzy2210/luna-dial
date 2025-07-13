package biz

import (
	"context"
	"strings"
	"testing"
	"time"
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

	if usecase == nil {
		t.Fatal("NewTaskUsecase returned nil")
	}

	if usecase.repo != repo {
		t.Error("repo not set correctly")
	}
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
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
			Tags:  []string{"工作", "文档", "产品"},
			Icon:  "📝",
			Score: 80,
		}

		task, err := usecase.CreateTask(ctx, param)

		// 期望成功创建，但当前会失败
		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: CreateTask 应该成功创建，但得到错误: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 业务逻辑未实现: CreateTask 应该返回创建的任务对象")
		}

		// 验证返回的任务字段
		if task.Title != param.Title {
			t.Errorf("期望标题为 %s, 得到 %s", param.Title, task.Title)
		}

		if task.TaskType != param.Type {
			t.Errorf("期望类型为 %v, 得到 %v", param.Type, task.TaskType)
		}

		if task.Score != param.Score {
			t.Errorf("期望分数为 %d, 得到 %d", param.Score, task.Score)
		}

		if task.UserID != param.UserID {
			t.Errorf("期望用户ID为 %s, 得到 %s", param.UserID, task.UserID)
		}

		if task.Icon != param.Icon {
			t.Errorf("期望图标为 %s, 得到 %s", param.Icon, task.Icon)
		}

		if len(task.Tags) != len(param.Tags) {
			t.Errorf("期望标签数量为 %d, 得到 %d", len(param.Tags), len(task.Tags))
		}

		// 验证自动设置的字段
		if task.ID == "" {
			t.Error("期望生成非空的ID")
		}

		if task.IsCompleted {
			t.Error("新创建的任务应该是未完成状态")
		}

		if task.CreatedAt.IsZero() {
			t.Error("期望设置创建时间")
		}

		if task.UpdatedAt.IsZero() {
			t.Error("期望设置更新时间")
		}
	})

	t.Run("成功创建周任务", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "完成项目里程碑",
			Type:   PeriodWeek,
			Period: Period{
				Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 19, 23, 59, 59, 0, time.UTC),
			},
			Tags:  []string{"项目", "里程碑"},
			Icon:  "🎯",
			Score: 200,
		}

		task, err := usecase.CreateTask(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回创建的周任务")
		}

		if task.TaskType != PeriodWeek {
			t.Errorf("期望任务类型为 PeriodWeek, 得到 %v", task.TaskType)
		}
	})

	t.Run("成功创建子任务", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "子任务：设计UI界面",
			Type:   PeriodDay,
			Period: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
			Tags:     []string{"设计", "UI"},
			Icon:     "🎨",
			Score:    50,
			ParentID: "parent-task-123", // 父任务ID
		}

		task, err := usecase.CreateTask(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回创建的子任务")
		}

		if task.ParentID != param.ParentID {
			t.Errorf("期望父任务ID为 %s, 得到 %s", param.ParentID, task.ParentID)
		}
	})

	t.Run("参数验证失败 - 空用户ID", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "", // 空用户ID
			Title:  "测试任务",
			Type:   PeriodDay,
			Score:  50,
		}

		task, err := usecase.CreateTask(ctx, param)

		if task != nil {
			t.Errorf("期望返回 nil, 得到 %+v", task)
		}

		if err == nil {
			t.Error("期望返回验证错误")
		}

		// TODO: 实现后应该返回具体的验证错误
		if err == ErrNoPermission {
			t.Log("当前返回 ErrNoPermission，实现后应该返回具体的验证错误")
		}
	})

	t.Run("参数验证失败 - 空标题", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "", // 空标题
			Type:   PeriodDay,
			Score:  50,
		}

		task, err := usecase.CreateTask(ctx, param)

		if task != nil {
			t.Errorf("期望返回 nil, 得到 %+v", task)
		}

		if err == nil {
			t.Error("期望返回验证错误")
		}
	})

	t.Run("参数验证失败 - 无效分数", func(t *testing.T) {
		param := CreateTaskParam{
			UserID: "user-123",
			Title:  "测试任务",
			Type:   PeriodDay,
			Score:  -10, // 负分数
		}

		task, err := usecase.CreateTask(ctx, param)

		if task != nil {
			t.Errorf("期望返回 nil, 得到 %+v", task)
		}

		if err == nil {
			t.Error("期望返回验证错误")
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: UpdateTask 应该成功更新，但得到错误: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回更新后的任务对象")
		}

		if task.Title != newTitle {
			t.Errorf("期望标题更新为 %s, 得到 %s", newTitle, task.Title)
		}

		// 验证更新时间被修改
		if task.UpdatedAt.IsZero() {
			t.Error("期望更新时间被设置")
		}
	})

	t.Run("成功更新任务完成状态", func(t *testing.T) {
		completed := true
		param := UpdateTaskParam{
			TaskID:      "task-123",
			UserID:      "user-123",
			IsCompleted: &completed,
		}

		task, err := usecase.UpdateTask(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回更新后的任务对象")
		}

		if !task.IsCompleted {
			t.Error("期望任务状态更新为已完成")
		}
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回更新后的任务对象")
		}

		if task.Score != newScore {
			t.Errorf("期望分数更新为 %d, 得到 %d", newScore, task.Score)
		}

		if len(task.Tags) != len(newTags) {
			t.Errorf("期望标签数量为 %d, 得到 %d", len(newTags), len(task.Tags))
		}
	})

	t.Run("权限验证失败 - 不同用户", func(t *testing.T) {
		newTitle := "恶意更新"
		param := UpdateTaskParam{
			TaskID: "task-123",
			UserID: "other-user", // 不同的用户ID
			Title:  &newTitle,
		}

		task, err := usecase.UpdateTask(ctx, param)

		if task != nil {
			t.Errorf("期望返回 nil, 得到 %+v", task)
		}

		if err == nil {
			t.Error("期望返回权限错误")
		}
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: DeleteTask 应该成功删除，但得到错误: %v", err)
		}
	})

	t.Run("权限验证失败", func(t *testing.T) {
		param := DeleteTaskParam{
			TaskID: "task-123",
			UserID: "other-user",
		}

		err := usecase.DeleteTask(ctx, param)

		if err == nil {
			t.Error("期望返回权限错误")
		}
	})

	t.Run("任务不存在", func(t *testing.T) {
		param := DeleteTaskParam{
			TaskID: "non-existent",
			UserID: "user-123",
		}

		err := usecase.DeleteTask(ctx, param)

		if err == nil {
			t.Error("期望返回不存在错误")
		}
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: SetTaskScore 应该成功设置，但得到错误: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回更新后的任务对象")
		}

		if task.Score != param.Score {
			t.Errorf("期望分数为 %d, 得到 %d", param.Score, task.Score)
		}
	})

	t.Run("无效分数", func(t *testing.T) {
		param := SetTaskScoreParam{
			TaskID: "task-123",
			UserID: "user-123",
			Score:  -50, // 负分数
		}

		task, err := usecase.SetTaskScore(ctx, param)

		if task != nil {
			t.Errorf("期望返回 nil, 得到 %+v", task)
		}

		if err == nil {
			t.Error("期望返回验证错误")
		}
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: CreateSubTask 应该成功创建，但得到错误: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回创建的子任务")
		}

		if task.ParentID != param.ParentID {
			t.Errorf("期望父任务ID为 %s, 得到 %s", param.ParentID, task.ParentID)
		}

		if task.Title != param.Title {
			t.Errorf("期望标题为 %s, 得到 %s", param.Title, task.Title)
		}
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

		if task != nil {
			t.Errorf("期望返回 nil, 得到 %+v", task)
		}

		if err == nil {
			t.Error("期望返回父任务不存在错误")
		}
	})
}

// 测试 AddTag 和 RemoveTag 方法
func TestTaskUsecase_TagOperations(t *testing.T) {
	usecase := createTestTaskUsecase()
	ctx := context.Background()

	t.Run("成功添加标签", func(t *testing.T) {
		param := AddTagParam{
			TaskID: "task-123",
			UserID: "user-123",
			Tag:    "新标签",
		}

		task, err := usecase.AddTag(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: AddTag 应该成功添加，但得到错误: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回更新后的任务对象")
		}

		// 验证标签被添加
		tagFound := false
		for _, tag := range task.Tags {
			if tag == param.Tag {
				tagFound = true
				break
			}
		}
		if !tagFound {
			t.Errorf("期望标签 %s 被添加", param.Tag)
		}
	})

	t.Run("成功移除标签", func(t *testing.T) {
		param := RemoveTagParam{
			TaskID: "task-123",
			UserID: "user-123",
			Tag:    "要移除的标签",
		}

		task, err := usecase.RemoveTag(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: RemoveTag 应该成功移除，但得到错误: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回更新后的任务对象")
		}

		// 验证标签被移除
		for _, tag := range task.Tags {
			if tag == param.Tag {
				t.Errorf("标签 %s 应该被移除", param.Tag)
			}
		}
	})

	t.Run("添加重复标签", func(t *testing.T) {
		param := AddTagParam{
			TaskID: "task-123",
			UserID: "user-123",
			Tag:    "已存在标签",
		}

		task, err := usecase.AddTag(ctx, param)

		// 实现后应该处理重复标签的情况
		if err == ErrNoPermission {
			t.Log("当前返回 ErrNoPermission，实现后需要处理重复标签")
		}

		if task != nil {
			// 验证不会添加重复标签
			tagCount := 0
			for _, tag := range task.Tags {
				if tag == param.Tag {
					tagCount++
				}
			}
			if tagCount > 1 {
				t.Errorf("不应该添加重复标签")
			}
		}
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: SetTaskIcon 应该成功设置，但得到错误: %v", err)
		}

		if task == nil {
			t.Fatal("❌ 应该返回更新后的任务对象")
		}

		if task.Icon != param.Icon {
			t.Errorf("期望图标为 %s, 得到 %s", param.Icon, task.Icon)
		}
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: ListTaskByPeriod 应该成功获取，但得到错误: %v", err)
		}

		if tasks == nil {
			t.Fatal("❌ 应该返回任务列表")
		}

		// 验证返回的任务都在指定时间范围内
		for _, task := range tasks {
			if task.UserID != param.UserID {
				t.Errorf("返回了其他用户的任务: %s", task.UserID)
			}

			// 验证任务时间在范围内
			if task.TimePeriod.Start.Before(param.Period.Start) ||
				task.TimePeriod.End.After(param.Period.End) {
				t.Errorf("任务时间超出范围: %v", task.TimePeriod)
			}
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: ListTaskTree 应该成功获取，但得到错误: %v", err)
		}

		if tasks == nil {
			t.Fatal("❌ 应该返回任务树列表")
		}

		// 验证任务树结构
		parentFound := false
		for _, task := range tasks {
			if task.ID == param.TaskID {
				parentFound = true
			}

			if task.UserID != param.UserID {
				t.Errorf("返回了其他用户的任务: %s", task.UserID)
			}
		}

		if !parentFound {
			t.Error("应该包含根任务")
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: ListTaskParentTree 应该成功获取，但得到错误: %v", err)
		}

		if tasks == nil {
			t.Fatal("❌ 应该返回父任务树列表")
		}

		// 验证返回的都是父级任务
		for _, task := range tasks {
			if task.UserID != param.UserID {
				t.Errorf("返回了其他用户的任务: %s", task.UserID)
			}
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

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: GetTaskStats 应该成功获取，但得到错误: %v", err)
		}

		if stats == nil {
			t.Fatal("❌ 应该返回统计数据")
		}

		// 期望返回12个月的统计数据
		expectedMonths := 12
		if len(stats) != expectedMonths {
			t.Errorf("期望 %d 个月的统计，得到 %d", expectedMonths, len(stats))
		}

		// 验证统计数据格式
		for _, stat := range stats {
			if stat.TaskCount < 0 {
				t.Errorf("任务数量不能为负数: %d", stat.TaskCount)
			}

			if stat.ScoreTotal < 0 {
				t.Errorf("总分不能为负数: %d", stat.ScoreTotal)
			}
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
		Tags:        []string{"测试", "任务"},
		Icon:        "📝",
		Score:       80,
		IsCompleted: false,
		ParentID:    "",
		UserID:      "user-123",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

	if task.IsCompleted {
		t.Error("期望任务为未完成状态")
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
