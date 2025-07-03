package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTaskType 测试任务类型枚举
func TestTaskType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TaskType
		valid    bool
	}{
		{"Valid Year", "year", TaskTypeYear, true},
		{"Valid Quarter", "quarter", TaskTypeQuarter, true},
		{"Valid Month", "month", TaskTypeMonth, true},
		{"Valid Week", "week", TaskTypeWeek, true},
		{"Valid Day", "day", TaskTypeDay, true},
		{"Valid Year Upper", "YEAR", TaskTypeYear, true},
		{"Valid Mixed Case", "Quarter", TaskTypeQuarter, true},
		{"Invalid Type", "invalid", "", false},
		{"Empty String", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试验证函数
			assert.Equal(t, tt.valid, IsValidTaskType(tt.input))

			// 测试解析函数
			if tt.valid {
				result, err := ParseTaskType(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				_, err := ParseTaskType(tt.input)
				assert.Error(t, err)
			}
		})
	}

	// 测试获取所有类型
	allTypes := GetAllTaskTypes()
	assert.Len(t, allTypes, 5)
	assert.Contains(t, allTypes, TaskTypeYear)
	assert.Contains(t, allTypes, TaskTypeDay)

	// 测试String方法
	assert.Equal(t, "year", TaskTypeYear.String())
}

// TestTaskStatus 测试任务状态枚举
func TestTaskStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TaskStatus
		valid    bool
	}{
		{"Valid Pending", "pending", TaskStatusPending, true},
		{"Valid InProgress", "in-progress", TaskStatusInProgress, true},
		{"Valid Completed", "completed", TaskStatusCompleted, true},
		{"Valid Upper Case", "PENDING", TaskStatusPending, true},
		{"Invalid Status", "invalid", "", false},
		{"Empty String", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试验证函数
			assert.Equal(t, tt.valid, IsValidTaskStatus(tt.input))

			// 测试解析函数
			if tt.valid {
				result, err := ParseTaskStatus(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				_, err := ParseTaskStatus(tt.input)
				assert.Error(t, err)
			}
		})
	}

	// 测试获取所有状态
	allStatuses := GetAllTaskStatuses()
	assert.Len(t, allStatuses, 3)
	assert.Contains(t, allStatuses, TaskStatusPending)
	assert.Contains(t, allStatuses, TaskStatusCompleted)

	// 测试String方法
	assert.Equal(t, "pending", TaskStatusPending.String())
}

// TestTimeScale 测试时间尺度枚举
func TestTimeScale(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TimeScale
		valid    bool
	}{
		{"Valid Day", "day", TimeScaleDay, true},
		{"Valid Week", "week", TimeScaleWeek, true},
		{"Valid Month", "month", TimeScaleMonth, true},
		{"Valid Quarter", "quarter", TimeScaleQuarter, true},
		{"Valid Year", "year", TimeScaleYear, true},
		{"Valid Upper Case", "DAY", TimeScaleDay, true},
		{"Invalid Scale", "invalid", "", false},
		{"Empty String", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试验证函数
			assert.Equal(t, tt.valid, IsValidTimeScale(tt.input))

			// 测试解析函数
			if tt.valid {
				result, err := ParseTimeScale(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				_, err := ParseTimeScale(tt.input)
				assert.Error(t, err)
			}
		})
	}

	// 测试获取所有时间尺度
	allScales := GetAllTimeScales()
	assert.Len(t, allScales, 5)
	assert.Contains(t, allScales, TimeScaleDay)
	assert.Contains(t, allScales, TimeScaleYear)

	// 测试String方法
	assert.Equal(t, "day", TimeScaleDay.String())
}

// TestEntryType 测试日志条目类型枚举
func TestEntryType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected EntryType
		valid    bool
	}{
		{"Valid PlanStart", "plan-start", EntryTypePlanStart, true},
		{"Valid Reflection", "reflection", EntryTypeReflection, true},
		{"Valid Summary", "summary", EntryTypeSummary, true},
		{"Valid Upper Case", "PLAN-START", EntryTypePlanStart, true},
		{"Invalid Type", "invalid", "", false},
		{"Empty String", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试验证函数
			assert.Equal(t, tt.valid, IsValidEntryType(tt.input))

			// 测试解析函数
			if tt.valid {
				result, err := ParseEntryType(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				_, err := ParseEntryType(tt.input)
				assert.Error(t, err)
			}
		})
	}

	// 测试获取所有条目类型
	allTypes := GetAllEntryTypes()
	assert.Len(t, allTypes, 3)
	assert.Contains(t, allTypes, EntryTypePlanStart)
	assert.Contains(t, allTypes, EntryTypeSummary)

	// 测试String方法
	assert.Equal(t, "plan-start", EntryTypePlanStart.String())
}

// TestTaskHierarchy 测试任务层级关系
func TestTaskHierarchy(t *testing.T) {
	tests := []struct {
		name       string
		parentType TaskType
		childType  TaskType
		valid      bool
	}{
		{"Year to Quarter", TaskTypeYear, TaskTypeQuarter, true},
		{"Quarter to Month", TaskTypeQuarter, TaskTypeMonth, true},
		{"Month to Week", TaskTypeMonth, TaskTypeWeek, true},
		{"Week to Day", TaskTypeWeek, TaskTypeDay, true},
		{"Year to Month", TaskTypeYear, TaskTypeMonth, false},
		{"Year to Day", TaskTypeYear, TaskTypeDay, false},
		{"Day to Week", TaskTypeDay, TaskTypeWeek, false},
		{"Day to anything", TaskTypeDay, TaskTypeDay, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidTaskHierarchy(tt.parentType, tt.childType)
			assert.Equal(t, tt.valid, result)
		})
	}

	// 测试获取允许的子类型
	yearChildren := GetAllowedChildTypes(TaskTypeYear)
	assert.Len(t, yearChildren, 1)
	assert.Contains(t, yearChildren, TaskTypeQuarter)

	dayChildren := GetAllowedChildTypes(TaskTypeDay)
	assert.Len(t, dayChildren, 0)
}

// TestParseTimeReference 测试时间引用解析
func TestParseTimeReference(t *testing.T) {
	tests := []struct {
		name      string
		timeRef   string
		scale     TimeScale
		wantError bool
	}{
		{
			name:      "Valid quarter",
			timeRef:   "2024-Q4",
			scale:     TimeScaleQuarter,
			wantError: false,
		},
		{
			name:      "Valid month",
			timeRef:   "2025-07",
			scale:     TimeScaleMonth,
			wantError: false,
		},
		{
			name:      "Valid year",
			timeRef:   "2024",
			scale:     TimeScaleYear,
			wantError: false,
		},
		{
			name:      "Valid day",
			timeRef:   "2025-07-15",
			scale:     TimeScaleDay,
			wantError: false,
		},
		{
			name:      "Invalid quarter format",
			timeRef:   "2024-Q5",
			scale:     TimeScaleQuarter,
			wantError: true,
		},
		{
			name:      "Invalid month format",
			timeRef:   "2025-13",
			scale:     TimeScaleMonth,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTimeReference(tt.timeRef, tt.scale)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)

			// 验证时间范围是否合理
			assert.True(t, result.Start.Before(result.End) || result.Start.Equal(result.End))
		})
	}
}

// TestParseQuarterReference 测试季度解析
func TestParseQuarterReference(t *testing.T) {
	result, err := ParseTimeReference("2024-Q4", TimeScaleQuarter)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Q4 应该是 10月1日 到 12月31日
	assert.Equal(t, 2024, result.Start.Year())
	assert.Equal(t, 10, int(result.Start.Month()))
	assert.Equal(t, 1, result.Start.Day())

	assert.Equal(t, 2024, result.End.Year())
	assert.Equal(t, 12, int(result.End.Month()))
	assert.Equal(t, 31, result.End.Day())
}

// TestParseMonthReference 测试月份解析
func TestParseMonthReference(t *testing.T) {
	result, err := ParseTimeReference("2025-07", TimeScaleMonth)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 2025年7月应该是 7月1日 到 7月31日
	assert.Equal(t, 2025, result.Start.Year())
	assert.Equal(t, 7, int(result.Start.Month()))
	assert.Equal(t, 1, result.Start.Day())

	assert.Equal(t, 2025, result.End.Year())
	assert.Equal(t, 7, int(result.End.Month()))
	assert.Equal(t, 31, result.End.Day())
}

// TestGenerateTimeLabels 测试时间标签生成
func TestGenerateTimeLabels(t *testing.T) {
	// 测试日级别标签生成
	timeRange := &TimeRange{
		Start: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2025, 7, 5, 23, 59, 59, 0, time.UTC),
	}

	labels := GenerateTimeLabels(timeRange, TimeScaleDay)

	assert.Len(t, labels, 5)
	assert.Equal(t, "2025-07-01", labels[0])
	assert.Equal(t, "2025-07-05", labels[4])
}

// TestTimeScaleValidation 测试时间尺度验证
func TestTimeScaleValidation(t *testing.T) {
	validScales := []TimeScale{
		TimeScaleDay,
		TimeScaleWeek,
		TimeScaleMonth,
		TimeScaleQuarter,
		TimeScaleYear,
	}

	for _, scale := range validScales {
		assert.True(t, scale.IsValid(), "Scale %v should be valid", scale)
	}

	invalidScale := TimeScale("invalid")
	assert.False(t, invalidScale.IsValid(), "Invalid scale should not be valid")
}
