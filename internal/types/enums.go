package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TaskType 任务类型枚举
type TaskType string

const (
	TaskTypeYear    TaskType = "year"    // 年度任务
	TaskTypeQuarter TaskType = "quarter" // 季度任务
	TaskTypeMonth   TaskType = "month"   // 月度任务
	TaskTypeWeek    TaskType = "week"    // 周任务
	TaskTypeDay     TaskType = "day"     // 日任务
)

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"     // 待开始
	TaskStatusInProgress TaskStatus = "in-progress" // 进行中
	TaskStatusCompleted  TaskStatus = "completed"   // 已完成
)

// TimeScale 时间尺度枚举
type TimeScale string

const (
	TimeScaleDay     TimeScale = "day"     // 日
	TimeScaleWeek    TimeScale = "week"    // 周
	TimeScaleMonth   TimeScale = "month"   // 月
	TimeScaleQuarter TimeScale = "quarter" // 季
	TimeScaleYear    TimeScale = "year"    // 年
)

// EntryType 日志条目类型枚举
type EntryType string

const (
	EntryTypePlanStart  EntryType = "plan-start" // 开始计划
	EntryTypeReflection EntryType = "reflection" // 阶段反思
	EntryTypeSummary    EntryType = "summary"    // 结束总结
)

// TaskType 验证和转换函数

// IsValidTaskType 验证任务类型是否有效
func IsValidTaskType(t string) bool {
	taskType := TaskType(strings.ToLower(t))
	switch taskType {
	case TaskTypeYear, TaskTypeQuarter, TaskTypeMonth, TaskTypeWeek, TaskTypeDay:
		return true
	default:
		return false
	}
}

// ParseTaskType 解析任务类型字符串
func ParseTaskType(t string) (TaskType, error) {
	taskType := TaskType(strings.ToLower(t))
	if !IsValidTaskType(t) {
		return "", fmt.Errorf("invalid task type: %s", t)
	}
	return taskType, nil
}

// GetAllTaskTypes 获取所有任务类型
func GetAllTaskTypes() []TaskType {
	return []TaskType{
		TaskTypeYear,
		TaskTypeQuarter,
		TaskTypeMonth,
		TaskTypeWeek,
		TaskTypeDay,
	}
}

// String 实现Stringer接口
func (t TaskType) String() string {
	return string(t)
}

// IsValid 验证任务类型是否有效
func (t TaskType) IsValid() bool {
	return IsValidTaskType(string(t))
}

// TaskStatus 验证和转换函数

// IsValidTaskStatus 验证任务状态是否有效
func IsValidTaskStatus(s string) bool {
	status := TaskStatus(strings.ToLower(s))
	switch status {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted:
		return true
	default:
		return false
	}
}

// ParseTaskStatus 解析任务状态字符串
func ParseTaskStatus(s string) (TaskStatus, error) {
	status := TaskStatus(strings.ToLower(s))
	if !IsValidTaskStatus(s) {
		return "", fmt.Errorf("invalid task status: %s", s)
	}
	return status, nil
}

// GetAllTaskStatuses 获取所有任务状态
func GetAllTaskStatuses() []TaskStatus {
	return []TaskStatus{
		TaskStatusPending,
		TaskStatusInProgress,
		TaskStatusCompleted,
	}
}

// String 实现Stringer接口
func (s TaskStatus) String() string {
	return string(s)
}

// IsValid 验证任务状态是否有效
func (s TaskStatus) IsValid() bool {
	return IsValidTaskStatus(string(s))
}

// TimeScale 验证和转换函数

// IsValidTimeScale 验证时间尺度是否有效
func IsValidTimeScale(scale string) bool {
	validScales := []string{
		string(TimeScaleDay),
		string(TimeScaleWeek),
		string(TimeScaleMonth),
		string(TimeScaleQuarter),
		string(TimeScaleYear),
	}

	// 支持大小写不敏感
	scale = strings.ToLower(scale)
	for _, validScale := range validScales {
		if scale == validScale {
			return true
		}
	}
	return false
}

// ParseTimeScale 解析时间尺度字符串
func ParseTimeScale(scale string) (TimeScale, error) {
	scale = strings.ToLower(strings.TrimSpace(scale))

	switch scale {
	case string(TimeScaleDay):
		return TimeScaleDay, nil
	case string(TimeScaleWeek):
		return TimeScaleWeek, nil
	case string(TimeScaleMonth):
		return TimeScaleMonth, nil
	case string(TimeScaleQuarter):
		return TimeScaleQuarter, nil
	case string(TimeScaleYear):
		return TimeScaleYear, nil
	default:
		return "", fmt.Errorf("无效的时间尺度: %s", scale)
	}
}

// GetAllTimeScales 获取所有有效的时间尺度
func GetAllTimeScales() []TimeScale {
	return []TimeScale{
		TimeScaleDay,
		TimeScaleWeek,
		TimeScaleMonth,
		TimeScaleQuarter,
		TimeScaleYear,
	}
}

// String 实现Stringer接口
func (ts TimeScale) String() string {
	return string(ts)
}

// IsValid 验证时间尺度是否有效
func (ts TimeScale) IsValid() bool {
	return IsValidTimeScale(string(ts))
}

// EntryType 验证和转换函数

// IsValidEntryType 验证日志条目类型是否有效
func IsValidEntryType(e string) bool {
	entryType := EntryType(strings.ToLower(e))
	switch entryType {
	case EntryTypePlanStart, EntryTypeReflection, EntryTypeSummary:
		return true
	default:
		return false
	}
}

// ParseEntryType 解析日志条目类型字符串
func ParseEntryType(e string) (EntryType, error) {
	entryType := EntryType(strings.ToLower(e))
	if !IsValidEntryType(e) {
		return "", fmt.Errorf("invalid entry type: %s", e)
	}
	return entryType, nil
}

// GetAllEntryTypes 获取所有日志条目类型
func GetAllEntryTypes() []EntryType {
	return []EntryType{
		EntryTypePlanStart,
		EntryTypeReflection,
		EntryTypeSummary,
	}
}

// String 实现Stringer接口
func (e EntryType) String() string {
	return string(e)
}

// IsValid 验证条目类型是否有效
func (e EntryType) IsValid() bool {
	return IsValidEntryType(string(e))
}

// 层级关系验证函数

// IsValidTaskHierarchy 验证任务层级关系是否有效
// 例如：年度任务可以包含季度任务，季度任务可以包含月度任务，等等
func IsValidTaskHierarchy(parentType, childType TaskType) bool {
	hierarchyMap := map[TaskType][]TaskType{
		TaskTypeYear:    {TaskTypeQuarter},
		TaskTypeQuarter: {TaskTypeMonth},
		TaskTypeMonth:   {TaskTypeWeek},
		TaskTypeWeek:    {TaskTypeDay},
		TaskTypeDay:     {}, // 日任务不能有子任务
	}

	allowedChildren, exists := hierarchyMap[parentType]
	if !exists {
		return false
	}

	for _, allowedChild := range allowedChildren {
		if allowedChild == childType {
			return true
		}
	}
	return false
}

// GetAllowedChildTypes 获取指定任务类型允许的子任务类型
func GetAllowedChildTypes(parentType TaskType) []TaskType {
	hierarchyMap := map[TaskType][]TaskType{
		TaskTypeYear:    {TaskTypeQuarter},
		TaskTypeQuarter: {TaskTypeMonth},
		TaskTypeMonth:   {TaskTypeWeek},
		TaskTypeWeek:    {TaskTypeDay},
		TaskTypeDay:     {}, // 日任务不能有子任务
	}

	return hierarchyMap[parentType]
}

// 时间参考解析和范围计算函数

// TimeRange 表示时间范围
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ParseTimeReference 解析时间参考字符串
// 支持格式: "2024-Q4", "2025-07", "2024", "2025-W15", "2025-07-15"
func ParseTimeReference(timeRef string, scale TimeScale) (*TimeRange, error) {
	timeRef = strings.TrimSpace(timeRef)

	switch scale {
	case TimeScaleYear:
		return parseYearReference(timeRef)
	case TimeScaleQuarter:
		return parseQuarterReference(timeRef)
	case TimeScaleMonth:
		return parseMonthReference(timeRef)
	case TimeScaleWeek:
		return parseWeekReference(timeRef)
	case TimeScaleDay:
		return parseDayReference(timeRef)
	default:
		return nil, fmt.Errorf("不支持的时间尺度: %s", scale)
	}
}

// parseYearReference 解析年份引用 (e.g., "2024")
func parseYearReference(timeRef string) (*TimeRange, error) {
	year, err := strconv.Atoi(timeRef)
	if err != nil {
		return nil, fmt.Errorf("无效的年份格式: %s", timeRef)
	}

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)

	return &TimeRange{Start: start, End: end}, nil
}

// parseQuarterReference 解析季度引用 (e.g., "2024-Q4")
func parseQuarterReference(timeRef string) (*TimeRange, error) {
	re := regexp.MustCompile(`^(\d{4})-Q([1-4])$`)
	matches := re.FindStringSubmatch(timeRef)
	if len(matches) != 3 {
		return nil, fmt.Errorf("无效的季度格式: %s, 期望格式: YYYY-QN", timeRef)
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("无效的年份: %s", matches[1])
	}

	quarter, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("无效的季度: %s", matches[2])
	}

	startMonth := (quarter-1)*3 + 1
	start := time.Date(year, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 3, 0).Add(-time.Nanosecond)

	return &TimeRange{Start: start, End: end}, nil
}

// parseMonthReference 解析月份引用 (e.g., "2025-07")
func parseMonthReference(timeRef string) (*TimeRange, error) {
	re := regexp.MustCompile(`^(\d{4})-(\d{2})$`)
	matches := re.FindStringSubmatch(timeRef)
	if len(matches) != 3 {
		return nil, fmt.Errorf("无效的月份格式: %s, 期望格式: YYYY-MM", timeRef)
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("无效的年份: %s", matches[1])
	}

	month, err := strconv.Atoi(matches[2])
	if err != nil || month < 1 || month > 12 {
		return nil, fmt.Errorf("无效的月份: %s", matches[2])
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return &TimeRange{Start: start, End: end}, nil
}

// parseWeekReference 解析周引用 (e.g., "2025-W15")
func parseWeekReference(timeRef string) (*TimeRange, error) {
	re := regexp.MustCompile(`^(\d{4})-W(\d{1,2})$`)
	matches := re.FindStringSubmatch(timeRef)
	if len(matches) != 3 {
		return nil, fmt.Errorf("无效的周格式: %s, 期望格式: YYYY-WN", timeRef)
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("无效的年份: %s", matches[1])
	}

	week, err := strconv.Atoi(matches[2])
	if err != nil || week < 1 || week > 53 {
		return nil, fmt.Errorf("无效的周数: %s", matches[2])
	}

	// 计算该年第一周的开始日期
	jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// 找到第一个周一
	weekday := jan1.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysToFirstMonday := time.Monday - weekday
	if daysToFirstMonday > 0 {
		daysToFirstMonday -= 7
	}

	firstMonday := jan1.AddDate(0, 0, int(daysToFirstMonday))
	start := firstMonday.AddDate(0, 0, (week-1)*7)
	end := start.AddDate(0, 0, 7).Add(-time.Nanosecond)

	return &TimeRange{Start: start, End: end}, nil
}

// parseDayReference 解析日期引用 (e.g., "2025-07-15")
func parseDayReference(timeRef string) (*TimeRange, error) {
	re := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
	matches := re.FindStringSubmatch(timeRef)
	if len(matches) != 4 {
		return nil, fmt.Errorf("无效的日期格式: %s, 期望格式: YYYY-MM-DD", timeRef)
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("无效的年份: %s", matches[1])
	}

	month, err := strconv.Atoi(matches[2])
	if err != nil || month < 1 || month > 12 {
		return nil, fmt.Errorf("无效的月份: %s", matches[2])
	}

	day, err := strconv.Atoi(matches[3])
	if err != nil || day < 1 || day > 31 {
		return nil, fmt.Errorf("无效的日期: %s", matches[3])
	}

	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1).Add(-time.Nanosecond)

	return &TimeRange{Start: start, End: end}, nil
}

// GenerateTimeLabels 根据时间范围和尺度生成时间标签
func GenerateTimeLabels(timeRange *TimeRange, scale TimeScale) []string {
	var labels []string
	current := timeRange.Start

	switch scale {
	case TimeScaleDay:
		for !current.After(timeRange.End) {
			labels = append(labels, current.Format("2006-01-02"))
			current = current.AddDate(0, 0, 1)
		}
	case TimeScaleWeek:
		// 周日-周六分组
		for !current.After(timeRange.End) {
			// 找到本周的周日
			startOfWeek := current
			if current.Weekday() != time.Sunday {
				startOfWeek = current.AddDate(0, 0, int(time.Sunday-current.Weekday()))
			}
			endOfWeek := startOfWeek.AddDate(0, 0, 6)
			if endOfWeek.After(timeRange.End) {
				endOfWeek = timeRange.End
			}
			label := startOfWeek.Format("2006-01-02") + "~" + endOfWeek.Format("2006-01-02")
			labels = append(labels, label)
			current = endOfWeek.AddDate(0, 0, 1)
		}
	case TimeScaleMonth:
		for !current.After(timeRange.End) {
			labels = append(labels, current.Format("2006-01"))
			current = current.AddDate(0, 1, 0)
		}
	case TimeScaleQuarter:
		for !current.After(timeRange.End) {
			quarter := (current.Month()-1)/3 + 1
			labels = append(labels, fmt.Sprintf("%d-Q%d", current.Year(), quarter))
			current = current.AddDate(0, 3, 0)
		}
	case TimeScaleYear:
		for !current.After(timeRange.End) {
			labels = append(labels, current.Format("2006"))
			current = current.AddDate(1, 0, 0)
		}
	}

	return labels
}
