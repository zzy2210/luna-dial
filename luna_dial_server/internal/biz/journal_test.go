package biz

import (
	"context"
	"strings"
	"testing"
	"time"
)

// 创建测试用的 JournalUsecase 实例
func createTestJournalUsecase() *JournalUsecase {
	repo := &mockJournalRepo{}
	return NewJournalUsecase(repo)
}

// 测试 NewJournalUsecase 构造函数
func TestNewJournalUsecase(t *testing.T) {
	repo := &mockJournalRepo{}
	usecase := NewJournalUsecase(repo)

	if usecase == nil {
		t.Fatal("NewJournalUsecase returned nil")
	}

	if usecase.repo != repo {
		t.Error("repo not set correctly")
	}
}

// 测试 CreateJournal 方法
func TestJournalUsecase_CreateJournal(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功创建日报", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       "今日工作总结",
			Content:     "今天完成了任务A和任务B，遇到了问题C并解决了。",
			JournalType: PeriodDay,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
			Icon: "📝",
		}

		journal, err := usecase.CreateJournal(ctx, param)

		// 期望成功创建，但当前会失败
		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: CreateJournal 应该成功创建，但得到错误: %v", err)
		}

		if journal == nil {
			t.Fatal("❌ 业务逻辑未实现: CreateJournal 应该返回创建的日志对象")
		}

		// 验证返回的日志字段
		if journal.Title != param.Title {
			t.Errorf("期望标题为 %s, 得到 %s", param.Title, journal.Title)
		}

		if journal.Content != param.Content {
			t.Errorf("期望内容为 %s, 得到 %s", param.Content, journal.Content)
		}

		if journal.JournalType != param.JournalType {
			t.Errorf("期望类型为 %v, 得到 %v", param.JournalType, journal.JournalType)
		}

		if journal.UserID != param.UserID {
			t.Errorf("期望用户ID为 %s, 得到 %s", param.UserID, journal.UserID)
		}

		if journal.Icon != param.Icon {
			t.Errorf("期望图标为 %s, 得到 %s", param.Icon, journal.Icon)
		}

		// 验证自动设置的字段
		if journal.ID == "" {
			t.Error("期望生成非空的ID")
		}

		if journal.CreatedAt.IsZero() {
			t.Error("期望设置创建时间")
		}

		if journal.UpdatedAt.IsZero() {
			t.Error("期望设置更新时间")
		}
	})

	t.Run("成功创建周报", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       "第3周工作总结",
			Content:     "本周完成了项目里程碑，团队协作效果良好。",
			JournalType: PeriodWeek,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 19, 23, 59, 59, 0, time.UTC),
			},
			Icon: "📊",
		}

		journal, err := usecase.CreateJournal(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if journal == nil {
			t.Fatal("❌ 应该返回创建的周报")
		}

		if journal.JournalType != PeriodWeek {
			t.Errorf("期望日志类型为 PeriodWeek, 得到 %v", journal.JournalType)
		}
	})

	t.Run("参数验证失败 - 空用户ID", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      "", // 空用户ID
			Title:       "测试日志",
			Content:     "测试内容",
			JournalType: PeriodDay,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
		}

		journal, err := usecase.CreateJournal(ctx, param)

		if journal != nil {
			t.Errorf("期望返回 nil, 得到 %+v", journal)
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
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       "", // 空标题
			Content:     "测试内容",
			JournalType: PeriodDay,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
			},
		}

		journal, err := usecase.CreateJournal(ctx, param)

		if journal != nil {
			t.Errorf("期望返回 nil, 得到 %+v", journal)
		}

		if err == nil {
			t.Error("期望返回验证错误")
		}
	})
}

// 测试 UpdateJournal 方法
func TestJournalUsecase_UpdateJournal(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功更新日志标题", func(t *testing.T) {
		newTitle := "更新后的标题"
		param := UpdateJournalParam{
			JournalID: "journal-123",
			UserID:    "user-123",
			Title:     &newTitle,
		}

		journal, err := usecase.UpdateJournal(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: UpdateJournal 应该成功更新，但得到错误: %v", err)
		}

		if journal == nil {
			t.Fatal("❌ 应该返回更新后的日志对象")
		}

		if journal.Title != newTitle {
			t.Errorf("期望标题更新为 %s, 得到 %s", newTitle, journal.Title)
		}

		// 验证更新时间被修改
		if journal.UpdatedAt.IsZero() {
			t.Error("期望更新时间被设置")
		}
	})

	t.Run("成功更新日志内容和类型", func(t *testing.T) {
		newContent := "更新后的内容"
		newType := PeriodWeek
		param := UpdateJournalParam{
			JournalID:   "journal-123",
			UserID:      "user-123",
			Content:     &newContent,
			JournalType: &newType,
		}

		journal, err := usecase.UpdateJournal(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if journal == nil {
			t.Fatal("❌ 应该返回更新后的日志对象")
		}

		if journal.Content != newContent {
			t.Errorf("期望内容更新为 %s, 得到 %s", newContent, journal.Content)
		}

		if journal.JournalType != newType {
			t.Errorf("期望类型更新为 %v, 得到 %v", newType, journal.JournalType)
		}
	})

	t.Run("权限验证失败 - 不同用户", func(t *testing.T) {
		newTitle := "恶意更新"
		param := UpdateJournalParam{
			JournalID: "journal-123",
			UserID:    "other-user", // 不同的用户ID
			Title:     &newTitle,
		}

		journal, err := usecase.UpdateJournal(ctx, param)

		if journal != nil {
			t.Errorf("期望返回 nil, 得到 %+v", journal)
		}

		if err == nil {
			t.Error("期望返回权限错误")
		}
	})
}

// 测试 DeleteJournal 方法
func TestJournalUsecase_DeleteJournal(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功删除日志", func(t *testing.T) {
		param := DeleteJournalParam{
			JournalID: "journal-123",
			UserID:    "user-123",
		}

		err := usecase.DeleteJournal(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: DeleteJournal 应该成功删除，但得到错误: %v", err)
		}
	})

	t.Run("权限验证失败", func(t *testing.T) {
		param := DeleteJournalParam{
			JournalID: "journal-123",
			UserID:    "other-user",
		}

		err := usecase.DeleteJournal(ctx, param)

		if err == nil {
			t.Error("期望返回权限错误")
		}
	})

	t.Run("日志不存在", func(t *testing.T) {
		param := DeleteJournalParam{
			JournalID: "non-existent",
			UserID:    "user-123",
		}

		err := usecase.DeleteJournal(ctx, param)

		if err == nil {
			t.Error("期望返回不存在错误")
		}
	})
}

// 测试 GetJournal 方法
func TestJournalUsecase_GetJournal(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功获取日志", func(t *testing.T) {
		param := GetJournalParam{
			JournalID: "journal-123",
			UserID:    "user-123",
		}

		journal, err := usecase.GetJournal(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: GetJournal 应该成功获取，但得到错误: %v", err)
		}

		if journal == nil {
			t.Fatal("❌ 应该返回日志对象")
		}

		if journal.ID != param.JournalID {
			t.Errorf("期望日志ID为 %s, 得到 %s", param.JournalID, journal.ID)
		}

		if journal.UserID != param.UserID {
			t.Errorf("期望用户ID为 %s, 得到 %s", param.UserID, journal.UserID)
		}
	})

	t.Run("日志不存在", func(t *testing.T) {
		param := GetJournalParam{
			JournalID: "non-existent",
			UserID:    "user-123",
		}

		journal, err := usecase.GetJournal(ctx, param)

		if journal != nil {
			t.Errorf("期望返回 nil, 得到 %+v", journal)
		}

		if err == nil {
			t.Error("期望返回不存在错误")
		}
	})
}

// 测试 ListJournalByPeriod 方法
func TestJournalUsecase_ListJournalByPeriod(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功获取月度日志列表", func(t *testing.T) {
		param := ListJournalByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodMonth,
		}

		journals, err := usecase.ListJournalByPeriod(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: ListJournalByPeriod 应该成功获取，但得到错误: %v", err)
		}

		if journals == nil {
			t.Fatal("❌ 应该返回日志列表")
		}

		// 验证返回的日志都在指定时间范围内
		for _, journal := range journals {
			if journal.UserID != param.UserID {
				t.Errorf("返回了其他用户的日志: %s", journal.UserID)
			}

			// 验证日志时间在范围内
			if journal.TimePeriod.Start.Before(param.Period.Start) ||
				journal.TimePeriod.End.After(param.Period.End) {
				t.Errorf("日志时间超出范围: %v", journal.TimePeriod)
			}
		}
	})

	t.Run("成功获取周度日志列表", func(t *testing.T) {
		param := ListJournalByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 19, 23, 59, 59, 0, time.UTC),
			},
			GroupBy: PeriodWeek,
		}

		journals, err := usecase.ListJournalByPeriod(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if journals == nil {
			t.Fatal("❌ 应该返回日志列表")
		}

		// 验证返回的日志类型
		for _, journal := range journals {
			if journal.JournalType != PeriodWeek && journal.JournalType != PeriodDay {
				t.Errorf("期望周报或日报，得到 %v", journal.JournalType)
			}
		}
	})
}

// 测试 ListAllJournals 方法
func TestJournalUsecase_ListAllJournals(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功获取分页日志列表", func(t *testing.T) {
		param := ListAllJournalsParam{
			UserID: "user-123",
			Pagination: PaginationParam{
				PageNum:  1,
				PageSize: 10,
			},
		}

		journals, err := usecase.ListAllJournals(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: ListAllJournals 应该成功获取，但得到错误: %v", err)
		}

		if journals == nil {
			t.Fatal("❌ 应该返回日志列表")
		}

		// 验证分页大小
		if len(journals) > param.Pagination.PageSize {
			t.Errorf("返回数量超过分页大小: %d > %d", len(journals), param.Pagination.PageSize)
		}

		// 验证所有日志都属于指定用户
		for _, journal := range journals {
			if journal.UserID != param.UserID {
				t.Errorf("返回了其他用户的日志: %s", journal.UserID)
			}
		}
	})

	t.Run("空结果分页", func(t *testing.T) {
		param := ListAllJournalsParam{
			UserID: "user-with-no-journals",
			Pagination: PaginationParam{
				PageNum:  1,
				PageSize: 10,
			},
		}

		journals, err := usecase.ListAllJournals(ctx, param)

		if err != nil {
			t.Errorf("❌ 业务逻辑未实现: %v", err)
		}

		if journals == nil {
			t.Fatal("❌ 应该返回空列表，而不是 nil")
		}

		if len(journals) != 0 {
			t.Errorf("期望返回空列表，得到 %d 个日志", len(journals))
		}
	})
}

// 测试结构体字段
func TestJournal_Fields(t *testing.T) {
	journal := Journal{
		ID:          "journal-123",
		Title:       "测试日志",
		Content:     "测试内容",
		JournalType: PeriodDay,
		TimePeriod: Period{
			Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC),
		},
		Icon:      "📝",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    "user-123",
	}

	if journal.ID != "journal-123" {
		t.Errorf("期望ID为 'journal-123', 得到 %s", journal.ID)
	}

	if journal.Title != "测试日志" {
		t.Errorf("期望标题为 '测试日志', 得到 %s", journal.Title)
	}

	if journal.JournalType != PeriodDay {
		t.Errorf("期望类型为 PeriodDay, 得到 %v", journal.JournalType)
	}

	if journal.UserID != "user-123" {
		t.Errorf("期望用户ID为 'user-123', 得到 %s", journal.UserID)
	}
}

// 测试参数结构体
func TestCreateJournalParam_Fields(t *testing.T) {
	param := CreateJournalParam{
		UserID:      "user-123",
		Title:       "新日志",
		Content:     "新内容",
		JournalType: PeriodWeek,
		TimePeriod: Period{
			Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 19, 23, 59, 59, 0, time.UTC),
		},
		Icon: "📊",
	}

	if param.UserID != "user-123" {
		t.Errorf("期望用户ID为 'user-123', 得到 %s", param.UserID)
	}

	if param.JournalType != PeriodWeek {
		t.Errorf("期望类型为 PeriodWeek, 得到 %v", param.JournalType)
	}
}

func TestUpdateJournalParam_Fields(t *testing.T) {
	newTitle := "更新标题"
	newContent := "更新内容"

	param := UpdateJournalParam{
		JournalID: "journal-123",
		UserID:    "user-123",
		Title:     &newTitle,
		Content:   &newContent,
	}

	if param.JournalID != "journal-123" {
		t.Errorf("期望日志ID为 'journal-123', 得到 %s", param.JournalID)
	}

	if param.Title == nil || *param.Title != newTitle {
		t.Errorf("期望标题为 '%s', 得到 %v", newTitle, param.Title)
	}

	if param.Content == nil || *param.Content != newContent {
		t.Errorf("期望内容为 '%s', 得到 %v", newContent, param.Content)
	}
}

// 测试分页参数
func TestPaginationParam_Fields(t *testing.T) {
	param := PaginationParam{
		PageNum:  2,
		PageSize: 20,
	}

	if param.PageNum != 2 {
		t.Errorf("期望页码为 2, 得到 %d", param.PageNum)
	}

	if param.PageSize != 20 {
		t.Errorf("期望页大小为 20, 得到 %d", param.PageSize)
	}
}

// 边界测试
func TestJournalUsecase_EdgeCases(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("nil context", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       "测试",
			Content:     "测试内容",
			JournalType: PeriodDay,
		}

		// 使用 context.TODO() 而不是 nil
		_, err := usecase.CreateJournal(context.TODO(), param)

		// 当前实现返回 ErrNoPermission，实现后可能需要处理特殊 context
		if err == nil {
			t.Log("实现后需要考虑特殊 context 的处理")
		}
	})

	t.Run("极长标题", func(t *testing.T) {
		longTitle := strings.Repeat("很长的标题", 1000)
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       longTitle,
			Content:     "测试内容",
			JournalType: PeriodDay,
		}

		journal, err := usecase.CreateJournal(ctx, param)

		// 实现后应该有标题长度限制
		if err == ErrNoPermission {
			t.Log("当前返回 ErrNoPermission，实现后应该有标题长度验证")
		}

		if journal != nil && len(journal.Title) > 200 {
			t.Errorf("标题可能过长，需要限制长度")
		}
	})

	t.Run("特殊字符处理", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       "测试<script>alert('xss')</script>",
			Content:     "内容包含特殊字符: & < > \" '",
			JournalType: PeriodDay,
			Icon:        "🚀💡🎯",
		}

		journal, err := usecase.CreateJournal(ctx, param)

		if err == ErrNoPermission {
			t.Log("当前返回 ErrNoPermission，实现后需要处理特殊字符转义")
		}

		if journal != nil {
			// 验证特殊字符被正确处理
			if strings.Contains(journal.Title, "<script>") {
				t.Error("可能存在XSS风险，需要转义HTML标签")
			}
		}
	})
}
