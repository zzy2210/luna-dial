package biz

import (
	"context"
	"strings"
	"testing"
	"time"

	"luna_dial/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试用常量 - UUID格式（无连字符）
const (
	TestUserID123            = "550e8400e29b41d4a716446655440000"
	TestUserIDOther          = "550e8400e29b41d4a716446655440001"
	TestUserIDWithNoJournals = "550e8400e29b41d4a716446655440002"
	TestJournalID123         = "123e4567e89b12d3a456426614174000"
	TestJournalIDNonExistent = "123e4567e89b12d3a456426614174001"
	TestJournalID1           = "123e4567e89b12d3a456426614174002"
)

// Mock JournalRepo 实现用于测试
type mockJournalRepo struct{}

func (m *mockJournalRepo) CreateJournal(ctx context.Context, journal *Journal) error {
	return nil
}

func (m *mockJournalRepo) UpdateJournal(ctx context.Context, journal *Journal) error {
	return nil
}

func (m *mockJournalRepo) DeleteJournalWithAuth(ctx context.Context, journalID, userID string) error {
	// 模拟记录不存在
	if journalID == TestJournalIDNonExistent || journalID == "non-existent" {
		return model.ErrRecordNotFound
	}

	// 模拟权限检查失败：用户ID不匹配时也返回记录不存在
	// 这符合数据库 WHERE id = ? AND user_id = ? 的行为
	if userID == TestUserIDOther || userID == "other-user" {
		return model.ErrRecordNotFound
	}

	return nil
}

func (m *mockJournalRepo) GetJournalWithAuth(ctx context.Context, journalID, userID string) (*Journal, error) {
	// 同样的逻辑
	if journalID == TestJournalIDNonExistent || journalID == "non-existent" {
		return nil, model.ErrRecordNotFound
	}

	if userID == TestUserIDOther || userID == "other-user" {
		return nil, model.ErrRecordNotFound
	}

	// 返回正常的模拟数据
	return &Journal{
		ID:          journalID,
		Title:       "测试日志",
		Content:     "测试内容",
		JournalType: PeriodDay,
		UserID:      userID,
		TimePeriod: Period{
			Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *mockJournalRepo) ListJournals(ctx context.Context, userID string, periodStart, periodEnd time.Time, journalType int) ([]*Journal, error) {
	if userID == TestUserIDWithNoJournals {
		return []*Journal{}, nil
	}
	// 返回模拟的日志列表
	return []*Journal{
		{
			ID:          TestJournalID1,
			Title:       "日志1",
			Content:     "内容1",
			JournalType: PeriodDay,
			UserID:      userID,
			TimePeriod: Period{
				Start: periodStart,
				End:   periodEnd,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

func (m *mockJournalRepo) ListAllJournals(ctx context.Context, userID string, offset, limit int) ([]*Journal, error) {
	if userID == TestUserIDWithNoJournals {
		return []*Journal{}, nil
	}

	// 模拟分页数据
	allJournals := []*Journal{
		{
			ID:          TestJournalID1,
			Title:       "日志1",
			Content:     "内容1",
			JournalType: PeriodDay,
			UserID:      userID,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          TestJournalID123,
			Title:       "日志2",
			Content:     "内容2",
			JournalType: PeriodWeek,
			UserID:      userID,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 模拟分页逻辑
	start := offset
	end := offset + limit
	if start >= len(allJournals) {
		return []*Journal{}, nil
	}
	if end > len(allJournals) {
		end = len(allJournals)
	}

	return allJournals[start:end], nil
}

func (m *mockJournalRepo) ListJournalsWithPagination(ctx context.Context, userID string, page, pageSize int, journalType *int, periodStart, periodEnd *time.Time) ([]*Journal, int64, error) {
	if userID == TestUserIDWithNoJournals {
		return []*Journal{}, 0, nil
	}

	// 模拟完整的分页数据
	allJournals := []*Journal{
		{
			ID:          TestJournalID1,
			Title:       "日志1",
			Content:     "内容1",
			JournalType: PeriodDay,
			UserID:      userID,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          TestJournalID123,
			Title:       "日志2",
			Content:     "内容2",
			JournalType: PeriodWeek,
			UserID:      userID,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          "test-journal-3",
			Title:       "日志3",
			Content:     "内容3",
			JournalType: PeriodMonth,
			UserID:      userID,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 应用过滤条件
	filteredJournals := []*Journal{}
	for _, journal := range allJournals {
		// 如果指定了日志类型过滤
		if journalType != nil && int(journal.JournalType) != *journalType {
			continue
		}
		
		// 如果指定了时间范围过滤
		if periodStart != nil && journal.TimePeriod.Start.Before(*periodStart) {
			continue
		}
		if periodEnd != nil && journal.TimePeriod.End.After(*periodEnd) {
			continue
		}
		
		filteredJournals = append(filteredJournals, journal)
	}

	total := int64(len(filteredJournals))

	// 计算分页
	offset := (page - 1) * pageSize
	end := offset + pageSize

	if offset >= len(filteredJournals) {
		return []*Journal{}, total, nil
	}
	if end > len(filteredJournals) {
		end = len(filteredJournals)
	}

	return filteredJournals[offset:end], total, nil
}

// 创建测试用的 JournalUsecase 实例
func createTestJournalUsecase() *JournalUsecase {
	repo := &mockJournalRepo{}
	return NewJournalUsecase(repo)
}

// 测试 NewJournalUsecase 构造函数
func TestNewJournalUsecase(t *testing.T) {
	repo := &mockJournalRepo{}
	usecase := NewJournalUsecase(repo)

	require.NotNil(t, usecase, "NewJournalUsecase should not return nil")
	assert.Equal(t, repo, usecase.repo, "repo should be set correctly")
}

// 测试 CreateJournal 方法
func TestJournalUsecase_CreateJournal(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功创建日报", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      TestUserID123,
			Title:       "今日工作总结",
			Content:     "今天完成了任务A和任务B，遇到了问题C并解决了。",
			JournalType: PeriodDay,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			},
			Icon: "📝",
		}

		journal, err := usecase.CreateJournal(ctx, param)

		// ❌ TDD: 期望成功创建，当前业务逻辑未实现会失败
		require.NoError(t, err, "CreateJournal should succeed")
		require.NotNil(t, journal, "CreateJournal should return created journal object")

		// 验证返回的日志字段
		assert.Equal(t, param.Title, journal.Title, "title should match")
		assert.Equal(t, param.Content, journal.Content, "content should match")
		assert.Equal(t, param.JournalType, journal.JournalType, "journal type should match")
		assert.Equal(t, param.UserID, journal.UserID, "user ID should match")
		assert.Equal(t, param.Icon, journal.Icon, "icon should match")

		// 验证自动设置的字段
		assert.NotEmpty(t, journal.ID, "ID should be generated")
		assert.False(t, journal.CreatedAt.IsZero(), "created time should be set")
		assert.False(t, journal.UpdatedAt.IsZero(), "updated time should be set")
	})

	t.Run("成功创建周报", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      TestUserID123,
			Title:       "第3周工作总结",
			Content:     "本周完成了项目里程碑，团队协作效果良好。",
			JournalType: PeriodWeek,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			Icon: "📊",
		}

		journal, err := usecase.CreateJournal(ctx, param)

		// ❌ TDD: 期望成功创建，当前业务逻辑未实现会失败
		require.NoError(t, err, "CreateJournal should succeed for week journal")
		require.NotNil(t, journal, "should return created week journal")
		assert.Equal(t, PeriodWeek, journal.JournalType, "journal type should be PeriodWeek")
	})

	t.Run("参数验证失败 - 空用户ID", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      "", // 空用户ID
			Title:       "测试日志",
			Content:     "测试内容",
			JournalType: PeriodDay,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			},
		}

		journal, err := usecase.CreateJournal(ctx, param)

		// ✅ TDD: 明确期望的业务错误
		assert.Nil(t, journal, "should return nil journal for empty user ID")
		assert.Equal(t, ErrNoPermission, err, "should return ErrUserIDEmpty for empty user ID")
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

		// ✅ TDD: 明确期望的业务错误
		assert.Nil(t, journal, "should return nil journal for empty title")
		assert.Equal(t, ErrInvalidInput, err, "should return ErrInvalidInput for empty title")
	})
}

// 测试 UpdateJournal 方法
func TestJournalUsecase_UpdateJournal(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功更新日志标题", func(t *testing.T) {
		newTitle := "更新后的标题"
		param := UpdateJournalParam{
			JournalID: TestJournalID123,
			UserID:    TestUserID123,
			Title:     &newTitle,
		}

		journal, err := usecase.UpdateJournal(ctx, param)

		// ❌ TDD: 期望成功更新，当前业务逻辑未实现会失败
		require.NoError(t, err, "UpdateJournal should succeed")
		require.NotNil(t, journal, "should return updated journal object")
		assert.Equal(t, newTitle, journal.Title, "title should be updated")
		assert.False(t, journal.UpdatedAt.IsZero(), "updated time should be set")
	})

	t.Run("成功更新日志内容和时间", func(t *testing.T) {
		newContent := "更新后的内容"
		newPeriod := Period{
			Start: time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC),
		}
		param := UpdateJournalParam{
			JournalID:  TestJournalID123,
			UserID:     TestUserID123,
			Content:    &newContent,
			TimePeriod: &newPeriod,
		}

		journal, err := usecase.UpdateJournal(ctx, param)

		// ❌ TDD: 期望成功更新，当前业务逻辑未实现会失败
		require.NoError(t, err, "UpdateJournal should succeed")
		require.NotNil(t, journal, "should return updated journal object")
		assert.Equal(t, newContent, journal.Content, "content should be updated")
		assert.Equal(t, newPeriod, journal.TimePeriod, "time period should be updated")
	})

	t.Run("权限验证失败 - 不同用户", func(t *testing.T) {
		newTitle := "恶意更新"
		param := UpdateJournalParam{
			JournalID: TestJournalID123,
			UserID:    TestUserIDOther, // 不同的用户ID
			Title:     &newTitle,
		}

		journal, err := usecase.UpdateJournal(ctx, param)

		// 期望返回业务层的"不存在"错误，而不是权限错误
		assert.Nil(t, journal, "should return nil journal for permission denied")
		assert.Equal(t, ErrJournalNotFound, err, "should return ErrJournalNotFound for different user")
	})
}

// 测试 DeleteJournal 方法
func TestJournalUsecase_DeleteJournal(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功删除日志", func(t *testing.T) {
		param := DeleteJournalParam{
			JournalID: TestJournalID123,
			UserID:    TestUserID123,
		}

		err := usecase.DeleteJournal(ctx, param)

		// ❌ TDD: 期望成功删除，当前业务逻辑未实现会失败
		assert.NoError(t, err, "DeleteJournal should succeed")
	})

	t.Run("权限验证失败", func(t *testing.T) {
		param := DeleteJournalParam{
			JournalID: "journal-123",
			UserID:    "other-user",
		}

		err := usecase.DeleteJournal(ctx, param)

		// 期望返回业务层的"不存在"错误，而不是权限错误
		assert.Equal(t, ErrJournalNotFound, err, "should return ErrJournalNotFound for unauthorized access")
	})

	t.Run("日志不存在", func(t *testing.T) {
		param := DeleteJournalParam{
			JournalID: "non-existent",
			UserID:    "user-123",
		}

		err := usecase.DeleteJournal(ctx, param)

		// ✅ TDD: 明确期望的不存在错误
		assert.Equal(t, ErrJournalNotFound, err, "should return ErrJournalNotFound for non-existent journal")
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

		// ❌ TDD: 期望成功获取，当前业务逻辑未实现会失败
		require.NoError(t, err, "GetJournal should succeed")
		require.NotNil(t, journal, "should return journal object")
		assert.Equal(t, param.JournalID, journal.ID, "journal ID should match")
		assert.Equal(t, param.UserID, journal.UserID, "user ID should match")
	})

	t.Run("权限验证失败", func(t *testing.T) {
		param := GetJournalParam{
			JournalID: "journal-123",
			UserID:    "other-user",
		}

		journal, err := usecase.GetJournal(ctx, param)

		assert.Nil(t, journal, "should return nil journal for unauthorized access")
		assert.Equal(t, ErrJournalNotFound, err, "should return ErrJournalNotFound for unauthorized access")
	})

	t.Run("日志不存在", func(t *testing.T) {
		param := GetJournalParam{
			JournalID: "non-existent",
			UserID:    "user-123",
		}

		journal, err := usecase.GetJournal(ctx, param)

		// ✅ TDD: 明确期望的不存在错误
		assert.Nil(t, journal, "should return nil journal for non-existent")
		assert.Equal(t, ErrJournalNotFound, err, "should return ErrJournalNotFound for non-existent journal")
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
				End:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			GroupBy: PeriodMonth,
		}

		journals, err := usecase.ListJournalByPeriod(ctx, param)

		// ❌ TDD: 期望成功获取，当前业务逻辑未实现会失败
		require.NoError(t, err, "ListJournalByPeriod should succeed")
		require.NotNil(t, journals, "should return journal list")

		// 验证返回的日志都在指定时间范围内
		for _, journal := range journals {
			assert.Equal(t, param.UserID, journal.UserID, "all journals should belong to specified user")
			// 验证日志时间在范围内
			assert.True(t, !journal.TimePeriod.Start.Before(param.Period.Start), "journal start time should be within range")
			assert.True(t, !journal.TimePeriod.End.After(param.Period.End), "journal end time should be within range")
		}
	})

	t.Run("成功获取周度日志列表", func(t *testing.T) {
		param := ListJournalByPeriodParam{
			UserID: "user-123",
			Period: Period{
				Start: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			GroupBy: PeriodWeek,
		}

		journals, err := usecase.ListJournalByPeriod(ctx, param)

		// ❌ TDD: 期望成功获取，当前业务逻辑未实现会失败
		require.NoError(t, err, "ListJournalByPeriod should succeed for week period")
		require.NotNil(t, journals, "should return journal list")

		// 验证返回的日志类型
		for _, journal := range journals {
			assert.Contains(t, []PeriodType{PeriodWeek, PeriodDay}, journal.JournalType, "should return week or day journals")
		}
	})
}

// 测试 ListAllJournals 方法
func TestJournalUsecase_ListAllJournals(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功获取分页日志列表", func(t *testing.T) {
		param := ListAllJournalsParam{
			UserID: TestUserID123,
			Pagination: PaginationParam{
				PageNum:  1,
				PageSize: 10,
			},
		}

		journals, err := usecase.ListAllJournals(ctx, param)

		// ✅ TDD: 期望成功获取
		require.NoError(t, err, "ListAllJournals should succeed")
		require.NotNil(t, journals, "should return journal list")

		// 验证分页大小
		assert.LessOrEqual(t, len(journals), param.Pagination.PageSize, "returned count should not exceed page size")

		// 验证所有日志都属于指定用户
		for _, journal := range journals {
			assert.Equal(t, param.UserID, journal.UserID, "all journals should belong to specified user")
		}
	})

	t.Run("空结果分页", func(t *testing.T) {
		param := ListAllJournalsParam{
			UserID: TestUserIDWithNoJournals,
			Pagination: PaginationParam{
				PageNum:  1,
				PageSize: 10,
			},
		}

		journals, err := usecase.ListAllJournals(ctx, param)

		// ✅ TDD: 期望成功获取空列表
		require.NoError(t, err, "ListAllJournals should succeed even with no results")
		require.NotNil(t, journals, "should return empty list, not nil")
		assert.Empty(t, journals, "should return empty list for user with no journals")
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

	assert.Equal(t, "journal-123", journal.ID, "ID should match")
	assert.Equal(t, "测试日志", journal.Title, "title should match")
	assert.Equal(t, PeriodDay, journal.JournalType, "journal type should match")
	assert.Equal(t, "user-123", journal.UserID, "user ID should match")
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

	assert.Equal(t, "user-123", param.UserID, "user ID should match")
	assert.Equal(t, PeriodWeek, param.JournalType, "journal type should match")
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

	assert.Equal(t, "journal-123", param.JournalID, "journal ID should match")
	require.NotNil(t, param.Title, "title pointer should not be nil")
	assert.Equal(t, newTitle, *param.Title, "title should match")
	require.NotNil(t, param.Content, "content pointer should not be nil")
	assert.Equal(t, newContent, *param.Content, "content should match")
}

// 测试分页参数
func TestPaginationParam_Fields(t *testing.T) {
	param := PaginationParam{
		PageNum:  2,
		PageSize: 20,
	}

	assert.Equal(t, 2, param.PageNum, "page number should match")
	assert.Equal(t, 20, param.PageSize, "page size should match")
}

// 测试 ListJournalsWithPagination 方法
func TestJournalUsecase_ListJournalsWithPagination(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("成功获取分页日志列表", func(t *testing.T) {
		param := ListJournalsWithPaginationParam{
			UserID:   TestUserID123,
			Page:     1,
			PageSize: 10,
		}

		journals, total, err := usecase.ListJournalsWithPagination(ctx, param)

		// ✅ TDD: 期望成功获取
		require.NoError(t, err, "ListJournalsWithPagination should succeed")
		require.NotNil(t, journals, "should return journal list")

		// 验证分页大小
		assert.LessOrEqual(t, len(journals), param.PageSize, "returned count should not exceed page size")
		assert.GreaterOrEqual(t, total, int64(len(journals)), "total should be greater than or equal to returned count")

		// 验证所有日志都属于指定用户
		for _, journal := range journals {
			assert.Equal(t, param.UserID, journal.UserID, "all journals should belong to specified user")
		}
	})

	t.Run("按日志类型过滤", func(t *testing.T) {
		dayType := int(PeriodDay)
		param := ListJournalsWithPaginationParam{
			UserID:      TestUserID123,
			Page:        1,
			PageSize:    10,
			JournalType: &dayType,
		}

		journals, total, err := usecase.ListJournalsWithPagination(ctx, param)

		// ✅ TDD: 期望成功获取过滤后的结果
		require.NoError(t, err, "ListJournalsWithPagination with filter should succeed")
		require.NotNil(t, journals, "should return filtered journal list")
		assert.GreaterOrEqual(t, total, int64(0), "total should be non-negative")

		// 验证所有返回的日志都是指定类型
		for _, journal := range journals {
			assert.Equal(t, PeriodDay, journal.JournalType, "all journals should be of specified type")
		}
	})

	t.Run("按时间范围过滤", func(t *testing.T) {
		startTime := time.Date(2025, 1, 14, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2025, 1, 16, 23, 59, 59, 0, time.UTC)
		
		param := ListJournalsWithPaginationParam{
			UserID:      TestUserID123,
			Page:        1,
			PageSize:    10,
			PeriodStart: &startTime,
			PeriodEnd:   &endTime,
		}

		journals, total, err := usecase.ListJournalsWithPagination(ctx, param)

		// ✅ TDD: 期望成功获取时间范围内的结果
		require.NoError(t, err, "ListJournalsWithPagination with time filter should succeed")
		require.NotNil(t, journals, "should return time-filtered journal list")
		assert.GreaterOrEqual(t, total, int64(0), "total should be non-negative")

		// 验证所有返回的日志都在指定时间范围内
		for _, journal := range journals {
			assert.True(t, !journal.TimePeriod.Start.Before(startTime), "journal start time should be within range")
			assert.True(t, !journal.TimePeriod.End.After(endTime), "journal end time should be within range")
		}
	})

	t.Run("空结果分页", func(t *testing.T) {
		param := ListJournalsWithPaginationParam{
			UserID:   TestUserIDWithNoJournals,
			Page:     1,
			PageSize: 10,
		}

		journals, total, err := usecase.ListJournalsWithPagination(ctx, param)

		// ✅ TDD: 期望成功获取空列表
		require.NoError(t, err, "ListJournalsWithPagination should succeed even with no results")
		require.NotNil(t, journals, "should return empty list, not nil")
		assert.Empty(t, journals, "should return empty list for user with no journals")
		assert.Equal(t, int64(0), total, "total should be 0 for user with no journals")
	})

	t.Run("参数验证失败 - 空用户ID", func(t *testing.T) {
		param := ListJournalsWithPaginationParam{
			UserID:   "", // 空用户ID
			Page:     1,
			PageSize: 10,
		}

		journals, total, err := usecase.ListJournalsWithPagination(ctx, param)

		// ✅ TDD: 明确期望的业务错误
		assert.Nil(t, journals, "should return nil for empty user ID")
		assert.Equal(t, int64(0), total, "total should be 0 for invalid input")
		assert.Equal(t, ErrUserIDEmpty, err, "should return ErrUserIDEmpty for empty user ID")
	})

	t.Run("参数自动调整 - 无效分页参数", func(t *testing.T) {
		param := ListJournalsWithPaginationParam{
			UserID:   TestUserID123,
			Page:     -1,    // 无效页码
			PageSize: 0,     // 无效页大小
		}

		journals, total, err := usecase.ListJournalsWithPagination(ctx, param)

		// ✅ TDD: 期望自动调整参数后成功
		require.NoError(t, err, "ListJournalsWithPagination should succeed with auto-adjusted params")
		require.NotNil(t, journals, "should return journal list with default pagination")
		assert.GreaterOrEqual(t, total, int64(0), "total should be non-negative")
		// 验证使用了默认的分页参数（页码1，页大小20）
		assert.LessOrEqual(t, len(journals), 20, "should use default page size of 20")
	})
}

// 边界测试
func TestJournalUsecase_EdgeCases(t *testing.T) {
	usecase := createTestJournalUsecase()
	ctx := context.Background()

	t.Run("极长标题", func(t *testing.T) {
		longTitle := strings.Repeat("很长的标题", 1000)
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       longTitle,
			Content:     "测试内容",
			JournalType: PeriodDay,
		}

		journal, err := usecase.CreateJournal(ctx, param)

		// ✅ TDD: 明确期望标题长度验证错误（未来需要定义具体错误类型）
		assert.Nil(t, journal, "should return nil journal for extremely long title")
		assert.Error(t, err, "should return validation error for extremely long title")
		// TODO: 实现后应该定义具体的标题长度错误类型
	})

	t.Run("空内容验证", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       "标题",
			Content:     "", // 空内容
			JournalType: PeriodDay,
		}

		journal, err := usecase.CreateJournal(ctx, param)

		// ✅ TDD: 明确期望内容为空的业务错误
		assert.Nil(t, journal, "should return nil journal for empty content")
		assert.Equal(t, ErrJournalContentEmpty, err, "should return ErrJournalContentEmpty for empty content")
	})

	t.Run("无效时间范围", func(t *testing.T) {
		param := CreateJournalParam{
			UserID:      "user-123",
			Title:       "测试",
			Content:     "测试内容",
			JournalType: PeriodDay,
			TimePeriod: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 14, 0, 0, 0, 0, time.UTC), // 结束时间在开始时间之前
			},
		}

		journal, err := usecase.CreateJournal(ctx, param)

		// ✅ TDD: 明确期望时间范围验证错误
		assert.Nil(t, journal, "should return nil journal for invalid time period")
		assert.Equal(t, ErrJournalPeriodInvalid, err, "should return ErrJournalPeriodInvalid for invalid period")
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

		// ❌ TDD: 期望特殊字符被正确处理，当前业务逻辑未实现会失败
		// 实现后应该能成功创建，但需要转义特殊字符
		if err == nil && journal != nil {
			// 验证特殊字符被正确处理
			assert.NotContains(t, journal.Title, "<script>", "should escape HTML tags to prevent XSS")
			assert.NotContains(t, journal.Content, "<", "should escape HTML characters")
		}
		// TODO: 实现后需要定义特殊字符处理的具体规则
	})
}
